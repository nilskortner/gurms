package aiserving

type AiServingProperties struct {
	adminApi *AdminApiProperties   `bson:"adminApiProperties"`
	mongo    *MongoGroupProperties `bson:"mongoGroupProperties"`
	ocr      *OcrProperties        `bson:"ocrProperties"`
}
