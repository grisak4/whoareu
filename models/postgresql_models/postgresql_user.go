package postgresqlmodels

import "time"

type User struct {
	ID          uint
	Username    string
	Nickname    string
	Email       string
	Password    string
	Role_System string
	Created_At  time.Time
	Updated_At  time.Time
}
