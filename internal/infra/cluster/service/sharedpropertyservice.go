package service

import (
	"context"
	"errors"
	"fmt"
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service/config/entity/configproperty"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property"
	"gurms/internal/storage/mongogurms/operation/option"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var SHAREDPROPERTYLOGGER logger.Logger = factory.GetLogger("SharedProperty")

type SharedPropertyService struct {
	clusterId string
	nodeType  nodetype.NodeType

	propertiesManager       *property.GurmsPropertiesManager
	sharedClusterProperties *configproperty.SharedClusterProperties
	sharedConfigService     *SharedConfigService

	propertiesChangeListeners []func(*property.GurmsProperties)

	cancelSharedPropertyService context.CancelFunc
}

func NewSharedPropertyService(clusterId string, nodeType nodetype.NodeType,
	propertiesManager *property.GurmsPropertiesManager) *SharedPropertyService {

	return &SharedPropertyService{
		clusterId:                 clusterId,
		nodeType:                  nodeType,
		propertiesManager:         propertiesManager,
		propertiesChangeListeners: make([]func(*property.GurmsProperties), 0),
	}
}

func (s *SharedPropertyService) LazyInit(sharedConfigService *SharedConfigService) {
	s.sharedConfigService = sharedConfigService
}

func (s *SharedPropertyService) Start() {
	go func() {
		opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
		stream, err := s.sharedConfigService.Subscribe(configproperty.SHAREDCLUSTERPROPERTIESNAME, opts)
		if err != nil {
			SHAREDPROPERTYLOGGER.FatalWithError("error subscribing to changestream", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		s.cancelSharedPropertyService = cancel
		for stream.Next(ctx) {
			var streamEvent bson.M
			if err := stream.Decode(&streamEvent); err != nil {
				SHAREDPROPERTYLOGGER.ErrorWithMessage("Error decoding change stream event:", err)
				continue
			}
			var changedProperties *configproperty.SharedClusterProperties
			if err := stream.Decode(&changedProperties); err != nil {
				SHAREDPROPERTYLOGGER.ErrorWithMessage("Error decoding to sharedclusterproperties:", err)
				continue
			}
			var changeClusterId string
			if changedProperties != nil {
				changeClusterId = changedProperties.ClusterId
			} else {
				if strAssert, ok := streamEvent["_id"].(string); ok {
					changeClusterId = strAssert
				}
			}

			if changeClusterId == s.clusterId {
				if operationType, ok := streamEvent["operationType"].(string); ok {
					switch operationType {
					case INSERT, REPLACE, UPDATE:
						s.sharedClusterProperties = changedProperties
						s.notifyListeners(s.sharedClusterProperties.GurmsProperties)
					case INVALIDATE:
						SHAREDPROPERTYLOGGER.Warn(
							"the shared properties has been deleted in mongodb unexpectedly")
						_, err := s.initializeSharedProperties()
						if err != nil {
							SHAREDPROPERTYLOGGER.ErrorWithMessage(
								"caught an error while initializing the shared properties", err)
						}
					default:
					}
				}
			}
		}
	}()

	_, err := s.initializeSharedProperties()
	if err != nil {
		SHAREDPROPERTYLOGGER.FatalWithError("failed to initialize the shared properties", err)
	}
}

func (s *SharedPropertyService) addChangeListener(listener func(*property.GurmsProperties)) {
	s.propertiesChangeListeners = append(s.propertiesChangeListeners, listener)
}

func (s *SharedPropertyService) notifyListeners(properties *property.GurmsProperties) {
	for _, listener := range s.propertiesChangeListeners {
		go listener(properties)
	}
}

func (s *SharedPropertyService) initializeSharedProperties() (*configproperty.SharedClusterProperties, error) {
	SHAREDPROPERTYLOGGER.InfoWithArgs("fetching shared properties")
	localProperties := s.propertiesManager.LocalGurmsProperties
	clusterProperties := &configproperty.SharedClusterProperties{
		ClusterId:       s.clusterId,
		SchemaVersion:   property.SCHEMA_VERSION,
		GurmsProperties: localProperties,
		LastUpdateTime:  time.Now(),
	}
	switch s.nodeType {
	case nodetype.AI_SERVING:
		clusterProperties.SetGatewayProperties(nil)
		clusterProperties.SetServiceProperties(nil)
		clusterProperties.SetMocknodeProperties(nil)
	case nodetype.GATEWAY:
		clusterProperties.SetAiServingProperties(nil)
		clusterProperties.SetServiceProperties(nil)
		clusterProperties.SetMocknodeProperties(nil)
	case nodetype.SERVICE:
		clusterProperties.SetAiServingProperties(nil)
		clusterProperties.SetGatewayProperties(nil)
		clusterProperties.SetMocknodeProperties(nil)
	case nodetype.MOCK:
		clusterProperties.SetAiServingProperties(nil)
		clusterProperties.SetGatewayProperties(nil)
		clusterProperties.SetServiceProperties(nil)
	}
	result, err := s.findAndUpdatePropertiesByNodeType(clusterProperties)
	if err == nil && result == nil {
		err = s.sharedConfigService.mongoClient.Operations.Insert(clusterProperties)
	}
	if err != nil {
		var writeErr mongo.WriteException
		errors.As(err, &writeErr)
		if writeErr.HasErrorCode(11000) {
			result, err = s.findAndUpdatePropertiesByNodeType(clusterProperties)
		}
		if err != nil {
			SHAREDPROPERTYLOGGER.ErrorWithMessage("failed to fetch shared properties", err)
			return nil, err
		}
	}
	s.sharedClusterProperties = result
	SHAREDPROPERTYLOGGER.InfoWithArgs("fetched shared properties")
	return result, nil
}

func (s *SharedPropertyService) findAndUpdatePropertiesByNodeType(
	clusterProperties *configproperty.SharedClusterProperties) (*configproperty.SharedClusterProperties, error) {

	filter := option.NewFilter()
	filter.Eq(configproperty.ID, s.clusterId)

	result := s.sharedConfigService.FindOne(configproperty.SHAREDCLUSTERPROPERTIESNAME, filter)
	var properties *configproperty.SharedClusterProperties
	err := result.Decode(&properties)
	if err != nil {
		err = fmt.Errorf("could not decode properties", err)
		return nil, err
	}
	switch s.nodeType {
	case nodetype.AI_SERVING:
		if properties.AiServingProperties != nil {
			return properties, nil
		}
		filter.Eq(configproperty.AISERVINGPROPERTIES, nil)
		update := option.NewUpdate()
		update.Set(configproperty.AISERVINGPROPERTIES, clusterProperties.AiServingProperties)
		result, err := s.sharedConfigService.mongoClient.Operations.UpdateOne(
			configproperty.SHAREDCLUSTERPROPERTIESNAME, filter, update)
		if err != nil {
			err = fmt.Errorf("failed to update properties", err)
			return nil, err
		}
		if result.MatchedCount == 0 {
			err = fmt.Errorf("failed to update properties", err)
			return nil, err
		}
		properties.SetAiServingProperties(clusterProperties.AiServingProperties)
		return properties, nil
	case nodetype.GATEWAY:
		if properties.GatewayProperties != nil {
			return properties, nil
		}
		filter.Eq(configproperty.GATEWAYPROPERTIES, nil)
		update := option.NewUpdate()
		update.Set(configproperty.GATEWAYPROPERTIES, clusterProperties.GatewayProperties)
		result, err := s.sharedConfigService.mongoClient.Operations.UpdateOne(
			configproperty.SHAREDCLUSTERPROPERTIESNAME, filter, update)
		if err != nil {
			err = fmt.Errorf("failed to update properties", err)
			return nil, err
		}
		if result.MatchedCount == 0 {
			err = fmt.Errorf("failed to update properties", err)
			return nil, err
		}
		properties.SetGatewayProperties(clusterProperties.GatewayProperties)
		return properties, nil
	case nodetype.SERVICE:
		if properties.ServiceProperties != nil {
			return properties, nil
		}
		filter.Eq(configproperty.SERVICEPROPERTIES, nil)
		update := option.NewUpdate()
		update.Set(configproperty.SERVICEPROPERTIES, clusterProperties.ServiceProperties)
		result, err := s.sharedConfigService.mongoClient.Operations.UpdateOne(
			configproperty.SHAREDCLUSTERPROPERTIESNAME, filter, update)
		if err != nil {
			err = fmt.Errorf("failed to update properties", err)
			return nil, err
		}
		if result.MatchedCount == 0 {
			err = fmt.Errorf("failed to update properties", err)
			return nil, err
		}
		properties.SetServiceProperties(clusterProperties.ServiceProperties)
		return properties, nil
	case nodetype.MOCK:
		if properties.MockNodeProperties != nil {
			return properties, nil
		}
		filter.Eq(configproperty.MOCKNODEPROPERTIES, nil)
		update := option.NewUpdate()
		update.Set(configproperty.MOCKNODEPROPERTIES, clusterProperties.MockNodeProperties)
		result, err := s.sharedConfigService.mongoClient.Operations.UpdateOne(
			configproperty.SHAREDCLUSTERPROPERTIESNAME, filter, update)
		if err != nil {
			err = fmt.Errorf("failed to update properties", err)
			return nil, err
		}
		if result.MatchedCount == 0 {
			err = fmt.Errorf("failed to update properties", err)
			return nil, err
		}
		properties.SetMocknodeProperties(clusterProperties.MockNodeProperties)
		return properties, nil
	default:
		err = fmt.Errorf("unknown nodetype")
		return nil, err
	}
}
