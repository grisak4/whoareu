package mongodbmodels

import (
	"time"
)

type Message struct {
	ID        uint      `json:"id" bson:"_id"`
	UserID    uint      `json:"user_id" bson:"user_id"`
	ChatID    uint      `json:"chat_id" bson:"chat_id"`
	Content   string    `json:"message_content" bson:"message_content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
