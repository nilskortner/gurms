package service

type ServiceProperties struct {
	// Env
	adminApi         *AdminApiProperties           `bson:",inline"`
	clientApi        *ClientApiProperties          `bson:",inline"`
	elasticsearch    *TurmsElasticsearchProperties `bson:",inline"`
	fake             *FakeProperties               `bson:",inline"`
	mongo            *MongoGroupProperties         `bson:",inline"`
	pushNotification *PushNotificationProperties   `bson:",inline"`
	redis            *TurmsRedisProperties         `bson:",inline"`
	statistics       *StatisticsProperties         `bson:",inline"`

	// Business
	conference   *ConferenceProperties   `bson:",inline"`
	conversation *ConversationProperties `bson:",inline"`
	message      *MessageProperties      `bson:",inline"`
	group        *GroupProperties        `bson:",inline"`
	user         *UserProperties         `bson:",inline"`
	storage      *StorageProperties      `bson:",inline"`
	notification *NotificationProperties `bson:",inline"`
}
