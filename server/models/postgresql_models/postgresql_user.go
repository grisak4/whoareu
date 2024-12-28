package postgresqlmodels

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	ID         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserConf   UserConfig `json:"user_config" gorm:"type:jsonb"`
	RoleSystem string     `json:"role_system"`
	Created_At time.Time  `json:"created_at"`
	Updated_At time.Time  `json:"updated_at"`
}

type UserConfig struct {
	Info    infoSettings    `json:"info"`
	Account accountSettings `json:"account"`
}

type infoSettings struct {
	Username string `json:"username" gorm:"unique"`
	Nickname string `json:"nickname"`
}
type accountSettings struct {
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

// Реализация интерфейса Scanner
func (u *UserConfig) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, u)
}

// Реализация интерфейса Valuer
func (u UserConfig) Value() (driver.Value, error) {
	return json.Marshal(u)
}
