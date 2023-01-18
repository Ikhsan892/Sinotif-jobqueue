package utils

import "os"

func GetWorkingDirectoryContent(filePath string) string {
	dir := os.Getenv("APP_DIRECTORY")
	return dir + filePath
}
