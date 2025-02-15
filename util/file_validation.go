package util

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/sidiqPratomo/max-health-backend/appconstant"
	"github.com/google/uuid"
)

func ValidateFile(fileHeader multipart.FileHeader, prefixPath string, formatOptions []string, maxSize int64) (*string, *string, error) {
	fileName := fileHeader.Filename
	splitFileName := strings.Split(fileName, ".")
	format := splitFileName[len(splitFileName)-1]

	if !isFileValidFormat(formatOptions, format) {
		return nil, nil, errors.New(appconstant.MsgInvalidFileType)
	}

	if fileHeader.Size > maxSize {
		return nil, nil, errors.New(appconstant.MsgTooLargeFile)
	}

	filePath := fmt.Sprintf("%s%s", prefixPath, uuid.NewString())
	return &filePath, &format, nil
}

func isFileValidFormat(formatOptions []string, format string) bool {
	for _, opt := range formatOptions {
		if opt == format {
			return true
		}
	}
	return false
}
