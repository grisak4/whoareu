package mongodbmodels

import "time"

type Chat struct {
	ID               uint      `json:"id" bson:"_id"`
	Title            string    `json:"title" bson:"title"`
	Participants_IDs []uint    `json:"participants_ids" bson:"participants_ids"`
	Type             string    `json:"type" bson:"type"` // conversation, group
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
}
