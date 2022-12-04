package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ORKL/library/olm/utils"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

const CDN_BASE_URL = "https://pub-7cb8ac806c1b4c4383e585c474a24719.r2.dev/"

func inCDN(sha1 string) bool {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := client.Head(CDN_BASE_URL + sha1 + ".pdf")
	if err != nil {
		if os.IsTimeout(err) {
			log.Error().Err(err).Msgf("making CDN HEAD request [timeout]: %v", sha1)
			return false
		} else {
			log.Error().Err(err).Msgf("making CDN HEAD request: %v", sha1)
			return false
		}
	}
	return res.StatusCode == 200
}

func getLibraryEntryHash(sha1 string) (Report, error) {
	var rhr ReportHashResponse
	var rep Report

	baseURL := "https://orkl.eu/api/v1/library/entry/sha1/"
	resp, err := http.Get(baseURL + sha1)
	if err != nil {
		return rep, err
	}
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rep, err
	}

	err = json.Unmarshal(body, &rhr)
	if err != nil {
		return rep, err
	}

	if rhr.Status == "error" {
		return rep, fmt.Errorf("supplied hash led to an API error: %v", rhr.Message)
	}

	rep.Sha1Hash = rhr.Data.Sha1Hash
	rep.Title = rhr.Data.Title
	rep.Authors = rhr.Data.Authors
	rep.FileNames = rhr.Data.ReportNames
	rep.FileSize = int64(rhr.Data.FileSize)
	rep.ReferenceURLs = rhr.Data.References
	rep.Language = rhr.Data.Language

	rep.PdfCreationDate = rhr.Data.FileCreationDate
	rep.PdfModificationDate = rhr.Data.FileModificationDate
	rep.PublicationDate = rhr.Data.FileCreationDate.Format("2006-01-02")

	rep.CDN = rhr.Data.Files.Pdf

	return rep, nil
}

func getLibraryEntryID(uuid string) (Report, error) {
	var rhr ReportHashResponse
	var rep Report

	baseURL := "https://orkl.eu/api/v1/library/entry/"
	resp, err := http.Get(baseURL + uuid)
	if err != nil {
		return rep, err
	}
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rep, err
	}

	err = json.Unmarshal(body, &rhr)
	if err != nil {
		return rep, err
	}

	if rhr.Status == "error" {
		return rep, fmt.Errorf("supplied uuid led to an API error: %v", rhr.Message)
	}

	rep.Sha1Hash = rhr.Data.Sha1Hash
	rep.Title = rhr.Data.Title
	rep.Authors = rhr.Data.Authors
	rep.FileNames = rhr.Data.ReportNames
	rep.FileSize = int64(rhr.Data.FileSize)
	rep.ReferenceURLs = rhr.Data.References
	rep.Language = rhr.Data.Language

	rep.PdfCreationDate = rhr.Data.FileCreationDate
	rep.PdfModificationDate = rhr.Data.FileModificationDate
	rep.PublicationDate = rhr.Data.FileCreationDate.Format("2006-01-02")

	rep.CDN = rhr.Data.Files.Pdf

	return rep, nil
}

func updateReportFromAPI(org Report) (Report, error) {
	apiRep, err := getLibraryEntryHash(org.Sha1Hash)
	if err != nil {
		return org, err
	}

	if org.Title == "" && apiRep.Title != "" {
		org.Title = apiRep.Title
	}

	if org.Authors == "" && apiRep.Authors != "" {
		org.Authors = apiRep.Authors
	}

	if org.FileNames == nil && apiRep.FileNames != nil {
		org.FileNames = apiRep.FileNames
	}

	if org.ReferenceURLs == nil && apiRep.ReferenceURLs != nil {
		org.ReferenceURLs = apiRep.ReferenceURLs
	}

	if org.Language == "" && apiRep.Language != "" {
		org.Language = apiRep.Language
	}

	return org, nil

}

func handleHash(sha1 string, corpusPath string) error {
	report, err := getLibraryEntryHash(sha1)
	if err != nil {
		return err
	}

	yamlName := report.Sha1Hash + ".yaml"
	yamlPath := filepath.Join(corpusPath, yamlName)

	yamlExists := utils.CheckFileExists(yamlPath)

	if !yamlExists || OverWriteMetadata {
		yamlData, err := yaml.Marshal(&report)
		if err != nil {
			return err
		}
		err = os.WriteFile(yamlPath, yamlData, 0664)
		if err != nil {
			return err
		}
	} else {
		log.Info().Msgf("YAML file already exists - use overwrite-metadata if you are sure you want to replace it")
	}
	return nil
}

func handleUUID(uuid string, corpusPath string) error {
	report, err := getLibraryEntryID(uuid)
	if err != nil {
		return err
	}

	yamlName := report.Sha1Hash + ".yaml"
	yamlPath := filepath.Join(corpusPath, yamlName)

	yamlExists := utils.CheckFileExists(yamlPath)

	if !yamlExists || OverWriteMetadata {
		yamlData, err := yaml.Marshal(&report)
		if err != nil {
			return err
		}
		err = os.WriteFile(yamlPath, yamlData, 0664)
		if err != nil {
			return err
		}
	} else {
		log.Info().Msgf("YAML file already exists - use overwrite-metadata if you are sure you want to replace it")
	}
	return nil
}
