package aiserving

type AiServingProperties struct {
	adminApi *AdminApiProperties   `bson:",inline"`
	mongo    *MongoGroupProperties `bson:",inline"`
	ocr      *OcrProperties        `bson:",inline"`
}
