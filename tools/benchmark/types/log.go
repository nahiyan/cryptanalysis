package types

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	InstanceName string `gorm:"index,instance_name"`
	Message      string
	IsValid      bool `gorm:"is_valid"`

	SessionID int
	Session
}
