package functions

import (
	"encoding/json"
	"fmt"
	"github.com/adjust/rmq/v5"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"math/rand"
	"sinotif/configs"
	"sinotif/pkg/utils"
	"time"
)

type ReportSmallTalk struct {
	db          *gorm.DB
	workerIndex int
	logProcess  *utils.JobQueueLog
	config      *configs.Config
}

var (
	processName    = "Report Small Talk"
	fileNamePrefix = "Data Hasil Small Talk "
	sheetName      = "Small Talk"
	headers        = []string{"No", "Nama Murid", "Tanggal", "Subject", "Detail", "Status"}
)

type PayloadReportSmallTalk struct {
	UserId    uint  `json:"user_id"`
	StudentId *uint `json:"student_id"`
}

type ResultReportSmallTalk struct {
	Tanggal     time.Time `json:"tanggal"`
	Subject     string    `json:"subject"`
	Detail      string    `json:"detail"`
	Status      string    `json:"status"`
	StudentName string    `json:"student_name"`
}

func NewReportSmalltalk(db *gorm.DB, workerIndex int, config *configs.Config) *ReportSmallTalk {
	return &ReportSmallTalk{
		db:          db,
		config:      config,
		workerIndex: workerIndex,
		logProcess: &utils.JobQueueLog{
			ProcessName:    processName,
			ProcessPayload: "",
			ProcessStatus:  "CREATED",
			ProcessType:    utils.REPORT,
			ProcessResult:  utils.FILE,
			IssuedBy:       0,
		},
	}
}

func ParsePayload(param string) (PayloadReportSmallTalk, error) {
	var payload PayloadReportSmallTalk
	// parsing payload
	if err := json.Unmarshal([]byte(param), &payload); err != nil {
		utils.Error(processName, err)
		return PayloadReportSmallTalk{}, err
	}
	return payload, nil
}

func (reportSmallTalk *ReportSmallTalk) GenerateExcel(param string) error {
	var payload PayloadReportSmallTalk
	// parsing payload
	if err := json.Unmarshal([]byte(param), &payload); err != nil {
		utils.Error(processName, err)
		return err
	}

	return nil
}

func (reportSmallTalk *ReportSmallTalk) Consume(delivery rmq.Delivery) {
	utils.Info(processName, fmt.Sprintf("Executing Job Report Small Talk on Worker %d...", reportSmallTalk.workerIndex))

	payload, err := ParsePayload(delivery.Payload())
	// parsing payload
	if err != nil {
		errReject := delivery.Reject()
		if errReject != nil {
			utils.Error(processName, errReject)
		}
	}

	// create first log
	reportSmallTalk.logProcess.ProcessPayload = delivery.Payload()
	reportSmallTalk.logProcess.IssuedBy = payload.UserId
	utils.CreateJobQueueLog(reportSmallTalk.db, reportSmallTalk.logProcess)

	// start the process
	reportSmallTalk.logProcess.ProcessStatus = "PROCESSING"
	utils.UpdateJobQueueLog(reportSmallTalk.db, reportSmallTalk.logProcess)

	// 1. fetch the data
	var resultReportSmallTalk []ResultReportSmallTalk

	condition := ""
	if payload.StudentId != nil {
		condition = fmt.Sprintf("and small_talks.student_id = %d", payload.StudentId)
	}

	if errSelect := reportSmallTalk.db.Raw("SELECT " +
		"small_talks.created_at as tanggal," +
		"small_talks.subject as subject," +
		"CONCAT(\"experienced_by_students_at_school\", ' | ', \"teaching_spec_student_service_at_sinotif\", ' | ', \"student_aspirations_targets_this_month\", ' | ', \"support_action_spec_aspirations_targets_students\") as detail," +
		"small_talks.status, " +
		"students.student_name " +
		"FROM small_talks LEFT JOIN students on students.id = small_talks.student_id " +
		"WHERE small_talks.deleted_at is null " +
		fmt.Sprintf("%s", condition)).Scan(&resultReportSmallTalk).Error; errSelect != nil {
		reportSmallTalk.logProcess.ProcessStatus = "FAILED"
		reportSmallTalk.logProcess.ProcessPayload = errSelect.Error()
		utils.UpdateJobQueueLog(reportSmallTalk.db, reportSmallTalk.logProcess)
		errReject := delivery.Reject()
		if errReject != nil {
			utils.Error(processName, errReject)
		}
	}

	// 2. process to excel
	file := excelize.NewFile()
	defer func() {
		if errExcel := file.Close(); errExcel != nil {
			utils.Error(processName, errExcel)
			errReject := delivery.Reject()
			if errReject != nil {
				utils.Error(processName, errReject)
			}
		}
	}()

	// 3. Create a new sheet.
	index, errFile := file.NewSheet(sheetName)
	if errFile != nil {
		utils.Error(processName, errFile)
		errReject := delivery.Reject()
		if errReject != nil {
			utils.Error(processName, errReject)
		}
	}

	file.SetActiveSheet(index)

	// 4. append header
	for i, header := range headers {
		colIndex, _ := excelize.ColumnNumberToName(i + 1)
		file.SetCellStr(sheetName, colIndex+"1", header)
	}
	style, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	colStart, _ := excelize.ColumnNumberToName(1)
	colEnd, _ := excelize.ColumnNumberToName(len(headers))
	file.SetCellStyle(sheetName, colStart+"1", colEnd+"1", style)

	// 5. append value
	for i, talk := range resultReportSmallTalk {
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), i+1)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), talk.StudentName)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), talk.Tanggal)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), talk.Subject)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), talk.Detail)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), talk.Status)
	}

	// 6. set excel location
	rand.Seed(time.Now().UTC().UnixNano())
	randomInt := rand.Int()
	fileName := fmt.Sprintf("%s/%s%d.xlsx", reportSmallTalk.config.Report.OutputLocation, fileNamePrefix, randomInt)
	file.Path = fileName
	file.DeleteSheet("Sheet1")

	// 7. save excel
	errSave := file.Save()
	if errSave != nil {
		utils.Error(processName, errSave)
	}

	// 8. update status
	reportSmallTalk.logProcess.ProcessStatus = "SUCCESS"
	reportSmallTalk.logProcess.ProcessPayload = fmt.Sprintf("%s%d.xlsx", fileNamePrefix, randomInt)
	utils.UpdateJobQueueLog(reportSmallTalk.db, reportSmallTalk.logProcess)

	utils.Info(processName, fmt.Sprintf("File Success Generated,workerIndex : %d, path : %s", reportSmallTalk.workerIndex, fileName))

	errAck := delivery.Ack()
	if errAck != nil {
		utils.Error(processName, errAck)
	}
}
