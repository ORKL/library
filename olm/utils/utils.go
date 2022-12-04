package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func SHA1HashForFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	f.Close()

	return hex.EncodeToString(h.Sum(nil)), nil
}

func GetFileSize(filepath string) (int64, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func CheckMimeType(path string) (string, error) {
	mType, err := mimetype.DetectFile(path)
	if err != nil {
		return "", err
	}
	return mType.String(), nil
}

func CheckIsPDF(path string) bool {
	mType, _ := CheckMimeType(path)
	return strings.Contains(strings.ToLower(mType), "pdf")
}

func EnsurePathExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckFileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
