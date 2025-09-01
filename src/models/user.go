package models

import (
	"strconv"
	"time"
)

type TUser struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Nickname   string    `json:"nickname" gorm:"not null;"`
	Password   string    `json:"password" gorm:"not null;"`
	Mobile     string    `json:"mobile" gorm:"not null;"`
	Email      string    `json:"email" gorm:"not null;"`
	CreateTime time.Time `json:"create_time" gorm:"not null;default:CURRENT_TIMESTAMP;type:datetime;"`
	ModifyTime time.Time `json:"modify_time" gorm:"not null;default:CURRENT_TIMESTAMP;type:datetime;"`
	Deleted    int       `json:"deleted" gorm:"deleted"`
}

func (user TUser) GetUid() string {
	return strconv.Itoa(int(user.ID))
}
