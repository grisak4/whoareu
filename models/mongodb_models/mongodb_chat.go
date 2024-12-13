package mongodbmodels

type Chat struct {
	ID               uint   `json:"chat_id" bson:"chat_id"`
	Title            string `json:"chat_title" bson:"chat_title"`
	Participants_IDs []uint `json:"participants_ids" bson:"participants_ids"`
	Type             string `json:"chat_type" bson:"chat_type"` // conversation, group
}
