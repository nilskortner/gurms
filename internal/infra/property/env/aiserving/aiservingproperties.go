package aiserving

type AiServingProperties struct {
	adminApi *AdminApiProperties   `bson:"adminApiProperties"`
	mongo    *MongoGroupProperties `bson:"mongoGroupProperties"`
	ocr      *OcrProperties        `bson:"ocrProperties"`
}

// TODO:
func NewAiServingProperties() *AiServingProperties {
	return &AiServingProperties{
		adminApi: NewAdminApiProperties(),
		mongo:    NewMongoGroupProperties(),
		ocr:      NewOcrProperties(),
	}
}
