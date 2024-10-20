package gormModel

import (
	"github.com/google/uuid"
)

type Account struct {
	BaseModel `swaggerignore:"true"`

	Name string    `json:"name" binding:"required" example:"Jonh Doe" gorm:"type:varchar(255)"`
	UID  uuid.UUID `json:"uid" example:"123e4567-e89b-12d3-a456-426614174000" gorm:"type:uuid;unique"`
}
