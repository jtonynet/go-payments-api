package gormModel

import (
	"github.com/google/uuid"
)

type Merchant struct {
	BaseModel `swaggerignore:"true"`
	UID       uuid.UUID `json:"uid" binding:"required" example:"0b0364a1-4955-48b6-8c63-8a446b918682" gorm:"type:uuid;uniqueIndex"`
	Name      string    `json:"name" binding:"required" example:"UBER EATS   SAO PAULO BR" gorm:"type:varchar(255);uniqueIndex"`
	MccID     uint      `json:"mcc_id" binding:"required" example:"1"`

	MCC MCC `gorm:"foreignKey:MccID"`
}
