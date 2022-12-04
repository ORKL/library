package main

import (
	"time"
)

type Report struct {
	Sha1Hash            string    `yaml:"sha1_hash"`
	Title               string    `yaml:"title"`
	Authors             string    `yaml:"authors"`
	PdfCreationDate     time.Time `yaml:"pdf_creation_date"`
	PdfModificationDate time.Time `yaml:"pdf_modification_date"`
	PublicationDate     string    `yaml:"publication_date"`
	FileNames           []string  `yaml:"file_names"`
	FileSize            int64     `yaml:"file_size"`
	ReferenceURLs       []string  `yaml:"reference_urls"`
	Language            string    `yaml:"language"`
	CDN                 string    `yaml:"cdn,omitempty"`
}

type ReportHashResponse struct {
	Data struct {
		ID                   string      `json:"id"`
		CreatedAt            time.Time   `json:"created_at"`
		UpdatedAt            time.Time   `json:"updated_at"`
		DeletedAt            interface{} `json:"deleted_at"`
		Sha1Hash             string      `json:"sha1_hash"`
		Title                string      `json:"title"`
		Authors              string      `json:"authors"`
		FileCreationDate     time.Time   `json:"file_creation_date"`
		FileModificationDate time.Time   `json:"file_modification_date"`
		FileSize             int         `json:"file_size"`
		PlainText            string      `json:"plain_text"`
		Language             string      `json:"language"`
		Sources              []struct {
			ID          string      `json:"id"`
			CreatedAt   time.Time   `json:"created_at"`
			UpdatedAt   time.Time   `json:"updated_at"`
			DeletedAt   interface{} `json:"deleted_at"`
			Name        string      `json:"name"`
			URL         string      `json:"url"`
			Description string      `json:"description"`
			Reports     interface{} `json:"reports"`
		} `json:"sources"`
		References         []string      `json:"references"`
		ReportNames        []string      `json:"report_names"`
		ThreatActors       []interface{} `json:"threat_actors"`
		TsCreatedAt        interface{}   `json:"ts_created_at"`
		TsUpdatedAt        interface{}   `json:"ts_updated_at"`
		TsCreationDate     interface{}   `json:"ts_creation_date"`
		TsModificationDate interface{}   `json:"ts_modification_date"`
		Files              struct {
			Pdf  string `json:"pdf"`
			Text string `json:"text"`
			Img  string `json:"img"`
		} `json:"files"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
