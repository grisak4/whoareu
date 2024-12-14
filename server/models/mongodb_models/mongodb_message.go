package mongodbmodels

type Message struct {
	ID      uint   `json:"message_id" bson:"message_id"`
	UserID  uint   `json:"user_id" bson:"user_id"`
	ChatID  uint   `json:"chat_id" bson:"chat_id"`
	Content string `json:"message_content" bson:"message_content"`
}
