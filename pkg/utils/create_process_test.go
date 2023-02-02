package utils

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

const (
	processName   = "TEST_LOG"
	processType   = REPORT
	processResult = FILE
	processStatus = CREATED
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateLog(t *testing.T) {
	payloadLog := LogInfo{
		ProcessName:    processName,
		ProcessPayload: "test",
		ProcessStatus:  processStatus,
		ProcessType:    processType,
		ProcessResult:  processResult,
		IssuedBy:       1,
	}

	db, mock, errMock := sqlmock.New()
	if errMock != nil {
		assert.Error(t, errMock)
	}

	open, errGorm := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}))
	if errGorm != nil {
		assert.Error(t, errGorm)
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery("INSERT INTO \"log_infos\" (.+)$").
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			payloadLog.ProcessName,
			payloadLog.ProcessPayload,
			payloadLog.ProcessStatus,
			payloadLog.ProcessType,
			payloadLog.ProcessResult,
			payloadLog.IssuedBy,
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	log := NewJobQueueLog(open)
	log.CreateJobQueueLog(payloadLog)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateLog(t *testing.T) {
	payloadLog := LogInfo{
		Model: gorm.Model{
			ID: 1,
		},
		ProcessName:    processName,
		ProcessPayload: "test",
		ProcessStatus:  processStatus,
		ProcessType:    processType,
		ProcessResult:  processResult,
		IssuedBy:       1,
	}

	db, mock, errMock := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if errMock != nil {
		assert.Error(t, errMock)
	}

	open, errGorm := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}))
	if errGorm != nil {
		assert.Error(t, errGorm)
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(regexp.QuoteMeta(`UPDATE "log_infos" SET "updated_at"=$1,"process_name"=$2,"process_payload"=$3,"process_status"=$4,"process_type"=$5,"process_result"=$6,"issued_by"=$7 WHERE id = $8 AND "log_infos"."deleted_at" IS NULL AND "id" = $9`)).
		WithArgs(
			AnyTime{},
			payloadLog.ProcessName,
			payloadLog.ProcessPayload,
			payloadLog.ProcessStatus,
			payloadLog.ProcessType,
			payloadLog.ProcessResult,
			payloadLog.IssuedBy,
			payloadLog.ID,
			payloadLog.ID,
		).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	log := NewJobQueueLog(open)
	log.UpdateJobQueueLog(payloadLog)

	err := mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
