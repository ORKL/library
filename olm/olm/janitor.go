package main

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// checkCDN checks for every yaml file in the corpus if the pdf report is in the
// CDN and deletes the PDF copy in the git when true
func checkCDN(corpPath string) {
	var yamlFiles []string

	err := filepath.Walk(corpPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error().Err(err).Msg("error in filepath walk")
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("error in filepath walk")
	}

	for _, file := range yamlFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Error().Err(err).Msg("error reading yaml file")
			continue
		}

		var report Report

		err = yaml.Unmarshal(content, &report)
		if err != nil {
			log.Error().Err(err).Msg("error unmarshal yaml file")
			continue
		}

		if report.CDN == "" {
			if inCDN(report.Sha1Hash) {
				report.CDN = CDN_BASE_URL + report.Sha1Hash + ".pdf"
				err = os.Remove(filepath.Join(corpPath, report.Sha1Hash+".pdf"))
				if err != nil {
					log.Error().Err(err).Msg("error removing report file")
					continue
				}
				yamlData, err := yaml.Marshal(&report)
				if err != nil {
					log.Error().Err(err).Msg("error marshal report yaml")
					continue
				}
				err = os.WriteFile(file, yamlData, 0664)
				if err != nil {
					log.Error().Err(err).Msg("error writing report yaml")
					continue
				}
			}
		}
	}
}
