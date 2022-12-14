package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ORKL/library/olm/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var OverWriteMetadata = false

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	debug := flag.Bool("debug", false, "sets log level to debug")
	corpPath := flag.String("library", "corpus", "path to library corpus folder")
	repPath := flag.String("report", "", "path to single report to import")
	recPath := flag.String("recursive", "", "path to directory with reports to import")
	hashStr := flag.String("hash", "", "populate metadata yaml from ORKL API for sha1 hash")
	uuidStr := flag.String("uuid", "", "populate metadata yaml from ORKL API for uuid")
	workInt := flag.Int("work", 0, "get specified number work items from the API")
	omBool := flag.Bool("overwrite-metadata", false, "overwrite yaml metadata files (caution)")
	janitorBool := flag.Bool("janitor", false, "run janitor to perform corpus maintenance and cleanup tasks")
	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	OverWriteMetadata = *omBool

	ok := checkCorpusPath(*corpPath)
	if !ok {
		log.Fatal().Msgf("path to corpus probably wrong: %v", *corpPath)
	}

	if *hashStr != "" {
		err := handleHash(*hashStr, *corpPath)
		if err != nil {
			log.Error().Err(err).Msgf("error handling supplied hash")
		}
	}

	if *uuidStr != "" {
		err := handleUUID(*uuidStr, *corpPath)
		if err != nil {
			log.Error().Err(err).Msgf("error handling supplied uuid")
		}
	}

	if *workInt > 0 {
		items, err := getLibraryWork(*workInt)
		if err != nil {
			log.Error().Err(err).Msgf("error getting work from the API")
			return
		}
		for _, item := range items {
			err = handleUUID(item.ID, *corpPath)
			if err != nil {
				log.Error().Err(err).Msgf("error handling supplied uuid")
			}
			niceReason := strings.Join(item.Reasons, " | ")
			fmt.Printf("Got %v to fix for you [%v]\n", item.ID, niceReason)
		}
	}

	if *repPath != "" {
		report, err := processReport(*repPath)
		if err != nil {
			log.Error().Err(err)
		}

		err = handleReport(report, *repPath, *corpPath)
		if err != nil {
			log.Error().Err(err)
		}
	}

	if *recPath != "" {
		err := filepath.Walk(*recPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Error().Err(err).Msgf("error in directory walk: %v", path)
				return nil
			}

			if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".pdf" {

				log.Debug().Msgf("processing PDF file: %v", path)

				report, err := processReport(path)
				if err != nil {
					log.Error().Err(err).Msgf("error processing report: %v", path)
				}

				handleErr := handleReport(report, path, *corpPath)
				if handleErr != nil {
					log.Error().Err(handleErr)
					return handleErr
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal().Err(err).Msg("error processing reports directory")
		}
	}

	if *janitorBool {
		checkCDN(*corpPath)
	}

	utils.EXIFTOOL.Close()
}

func checkCorpusPath(path string) bool {
	if _, err := os.Stat(filepath.Join(path, ".orkl_library")); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
