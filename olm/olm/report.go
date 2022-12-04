package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ORKL/library/olm/utils"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func init() {
	err := utils.InitiateExiftool()
	if err != nil {
		panic(err)
	}
}

func processReport(path string) (Report, error) {
	var report Report

	if !utils.CheckIsPDF(path) {
		return report, fmt.Errorf("error file is no PDF: %v", path)
	}

	size, err := utils.GetFileSize(path)
	if err != nil {
		return report, fmt.Errorf("error getting file size: %v [%v]", path, err)
	}
	report.FileSize = size

	hash, err := utils.SHA1HashForFile(path)
	if err != nil {
		return report, fmt.Errorf("error getting file hash: %v [%v]", path, err)
	}
	report.Sha1Hash = hash

	if info, err := utils.GetPDFMetadata(path); err != nil {
		return report, fmt.Errorf("exiftool file info: %v", err)
	} else {
		if val, ok := info.Fields["CreateDate"]; ok {
			createStr := fmt.Sprintf("%v", val)
			cd, err := utils.ConvertExiftoolDate(createStr)
			if err != nil {
				log.Warn().Err(err).Msgf("error checking creation date for file: %v", path)
			}
			report.PdfCreationDate = cd
			report.PublicationDate = cd.Format("2006-01-02")
		}

		if val, ok := info.Fields["ModifyDate"]; ok {
			modifiedStr := fmt.Sprintf("%v", val)
			md, err := utils.ConvertExiftoolDate(modifiedStr)
			if err != nil {
				log.Warn().Err(err).Msgf("error checking modified date for file: %v", path)
			}
			report.PdfModificationDate = md
		}

		if val, ok := info.Fields["Title"]; ok {
			report.Title = strings.TrimSpace(fmt.Sprintf("%v", val))
		}
	}

	fileName := filepath.Base(path)

	if fileName != hash+".pdf" {
		report.FileNames = append(report.FileNames, fileName)
	}

	if inCDN(report.Sha1Hash) {
		report.CDN = CDN_BASE_URL + report.Sha1Hash + ".pdf"
		report, err = updateReportFromAPI(report)
		if err != nil {
			return report, err
		}
	}

	return report, nil
}

func handleReport(report Report, original string, corpusPath string) error {
	pdfName := report.Sha1Hash + ".pdf"
	pdfPath := filepath.Join(corpusPath, pdfName)
	yamlName := report.Sha1Hash + ".yaml"
	yamlPath := filepath.Join(corpusPath, yamlName)

	pdfExists := utils.CheckFileExists(pdfPath)
	yamlExists := utils.CheckFileExists(yamlPath)

	if !pdfExists && report.CDN == "" {
		// Open original file
		original, err := os.Open(original)
		if err != nil {
			return err
		}
		defer original.Close()

		new, err := os.Create(pdfPath)
		if err != nil {
			return err
		}
		defer new.Close()

		_, err = io.Copy(new, original)
		if err != nil {
			return err
		}
	}

	if !yamlExists || OverWriteMetadata {
		yamlData, err := yaml.Marshal(&report)
		if err != nil {
			return err
		}
		err = os.WriteFile(yamlPath, yamlData, 0664)
		if err != nil {
			return err
		}
	}

	return nil
}
