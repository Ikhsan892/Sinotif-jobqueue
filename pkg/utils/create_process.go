package utils

import (
	"gorm.io/gorm"
	"os"
)

type JobQueueLog struct {
	gorm.Model
	ProcessName    string `json:"process_name"`
	ProcessPayload string `json:"process_payload"`
	ProcessStatus  string `json:"process_status"`
	ProcessType    string `json:"process_type"`
	ProcessResult  string `json:"process_result"`
	IssuedBy       uint   `json:"issued_by"`
}

// CreateJobQueueLog create log process
func CreateJobQueueLog(db *gorm.DB, data *JobQueueLog) {
	tx := db.Begin()
	if errCreate := tx.Table(os.Getenv("TABLE_LOG")).Create(&data).Error; errCreate != nil {
		Error("CREATE_TABLE_LOG_PROCESS", errCreate)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

// UpdateJobQueueLog update log process
func UpdateJobQueueLog(db *gorm.DB, data *JobQueueLog) {
	if errUpdate := db.Table(os.Getenv("TABLE_LOG")).Where("id = ?", data.ID).Updates(&data).Error; errUpdate != nil {
		Error("CREATE_TABLE_LOG_PROCESS", errUpdate)
	}
}
