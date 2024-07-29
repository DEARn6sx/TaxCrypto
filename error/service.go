package errorCase

type ResponseError struct {
	Status  int    `json:"status" bson:"status"`
	Code    int    `json:"code" bson:"code"`
	LabelTH string `json:"label_th" bson:"label_th"`
	LabelEN string `json:"label_en" bson:"label_en"`
}
