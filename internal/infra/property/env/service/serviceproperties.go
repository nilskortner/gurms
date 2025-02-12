package service

import (
	"gurms/internal/infra/property/env/service/business/conference"
	"gurms/internal/infra/property/env/service/business/conversation"
	"gurms/internal/infra/property/env/service/business/group"
	"gurms/internal/infra/property/env/service/business/message"
	"gurms/internal/infra/property/env/service/business/notification"
	"gurms/internal/infra/property/env/service/business/storage"
	"gurms/internal/infra/property/env/service/business/user"
	"gurms/internal/infra/property/env/service/env"
	"gurms/internal/infra/property/env/service/env/adminapi"
	"gurms/internal/infra/property/env/service/env/clientapi"
	"gurms/internal/infra/property/env/service/env/database"
	"gurms/internal/infra/property/env/service/env/elasticsearch"
	"gurms/internal/infra/property/env/service/env/push"
	"gurms/internal/infra/property/env/service/env/redis"
)

type ServiceProperties struct {
	// Env
	adminApi         *adminapi.AdminApiProperties                `bson:"adminApiProperties"`
	clientApi        *clientapi.ClientApiProperties              `bson:"clientApiProperties"`
	elasticsearch    *elasticsearch.GurmsElasticsearchProperties `bson:"turmsElasticsearchProperties"`
	fake             *env.FakeProperties                         `bson:"fakeProperties"`
	mongo            *database.MongoProperties                   `bson:"mongoGroupProperties"`
	pushNotification *push.PushNotificationProperties            `bson:"pushNotificationProperties"`
	redis            *redis.GurmsRedisProperties                 `bson:"turmsRedisProperties"`
	statistics       *env.StatisticsProperties                   `bson:"statisticsProperties"`

	// Business
	conference   *conference.ConferenceProperties     `bson:"conferenceProperties"`
	conversation *conversation.ConversationProperties `bson:"conversationProperties"`
	message      *message.MessageProperties           `bson:"messageProperties"`
	group        *group.GroupProperties               `bson:"groupProperties"`
	user         *user.UserProperties                 `bson:"userProperties"`
	storage      *storage.StorageProperties           `bson:"storageProperties"`
	notification *notification.NotificationProperties `bson:"notificationProperties"`
}

func NewServiceProperties() *ServiceProperties {
	return &ServiceProperties{
		adminApi:         adminapi.NewAdminApiProperties(),
		clientApi:        clientapi.NewClientApiProperties(),
		elasticsearch:    elasticsearch.NewGurmsElasticsearchProperties(),
		fake:             env.NewFakeProperties(),
		mongo:            database.NewMongoProperties(),
		pushNotification: push.NewPushNotificationProperties(),
		redis:            redis.NewGurmsRedisProperties(),
		statistics:       env.NewStatisticsProperties(),
		conference:       conference.NewConferenceProperties(),
		conversation:     conversation.NewConversationProperties(),
		message:          message.NewMessageProperties(),
		group:            group.NewGroupProperties(),
		user:             user.NewUserProperties(),
		storage:          storage.NewStorageProperties(),
		notification:     notification.NewNotificationProperties(),
	}
}
