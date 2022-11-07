package types

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	FileName string `gorm:"index,file_name"`
	SlurmId  uint   `gorm:"slurm_id"`
}
