package files

import (
	"os"
	"regexp"

	"github.com/pkg/errors"
	"olympos.io/encoding/edn"
)

func ReadEDNFile[T any](filePath string) (*FileContent[T], error) {
	fileInfo, err := ReadFileInfo(filePath)
	if err != nil {
		return nil, err
	}

	rawContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	var content T
	err = edn.Unmarshal(rawContent, &content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal edn content")
	}

	return &FileContent[T]{
		FileInfo: *fileInfo,
		Content:  content,
	}, nil
}

func ReadEDNDirs[T any](dirNames []string, regex *regexp.Regexp) ([]FileContent[T], error) {
	return readDirs[FileContent[T]](dirNames, regex, ReadEDNFile[T])
}

func WriteEDNFile(filePath string, content interface{}) error {
	data, err := edn.Marshal(content)
	if err != nil {
		return errors.Wrap(err, "failed to marshal end")
	}
	return WriteFile(filePath, data)
}
