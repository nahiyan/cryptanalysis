package types

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	FileName string `gorm:"index,file_name"`
	IdSlurm  uint   `gorm:"id_slurm"`

	SessionID int
	Session
}
