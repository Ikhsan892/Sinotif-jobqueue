package configs

import "os"

// ReportConfig TODO : To Create Reporting data
type ReportConfig struct {
	OutputLocation string
}

func LoadReportConfig() ReportConfig {
	return ReportConfig{
		OutputLocation: os.Getenv("FILE_REPORT_LOCATION"),
	}
}
