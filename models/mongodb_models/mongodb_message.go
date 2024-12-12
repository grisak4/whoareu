package mongodbmodels

type MongoDB_Message struct {
	UserID  uint   `bson:"user_id"`
	ChatID  uint   `bson:"chat_id"`
	Content string `bson:"content"`
}
