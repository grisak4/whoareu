package postgresqlmodels

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username   string    `json:"username" gorm:"not null;unique"`
	Nickname   string    `json:"nickname"`
	Email      string    `json:"email" gorm:"not null;unique"`
	Password   string    `json:"password" gorm:"not null;unique"`
	RoleSystem string    `json:"role_system"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}
