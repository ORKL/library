package utils

import (
	"fmt"
	"time"

	"github.com/barasher/go-exiftool"
)

var EXIFTOOL *exiftool.Exiftool

func InitiateExiftool() error {
	var err error
	buf := make([]byte, 128*1000)
	EXIFTOOL, err = exiftool.NewExiftool(
		exiftool.Buffer(buf, 64*1000),
		exiftool.DateFormant("%Y-%m-%dT%H:%M:%SZ"))
	if err != nil {
		return err
	}
	return nil
}

func ConvertExiftoolDate(dateString string) (time.Time, error) {
	dt, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return time.Time{}, err
	}
	return dt, nil
}

func GetPDFMetadata(path string) (exiftool.FileMetadata, error) {
	fileInfos := EXIFTOOL.ExtractMetadata(path)
	if fileInfos[0].Err != nil {
		return exiftool.FileMetadata{}, fmt.Errorf("Error when using exiftool: %v\n", fileInfos[0].Err)

	}
	return fileInfos[0], nil
}
