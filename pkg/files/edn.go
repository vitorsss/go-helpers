package files

import (
	"bufio"
	"bytes"
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
	return readDirs(dirNames, regex, ReadEDNFile[T])
}

func WriteEDNFile(filePath string, content interface{}) error {
	data, err := edn.Marshal(content)
	if err != nil {
		return errors.Wrap(err, "failed to marshal end")
	}
	return WriteFile(filePath, data)
}

func ReadNDEDNFile[T any](filePath string) (*FileContent[[]T], error) {
	fileInfo, err := ReadFileInfo(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)

	var content []T

	for scanner.Scan() {
		var contentLine T
		err = edn.Unmarshal([]byte(scanner.Text()), &contentLine)
		if err != nil {
			return nil, err
		}
		content = append(content, contentLine)
	}

	return &FileContent[[]T]{
		FileInfo: *fileInfo,
		Content:  content,
	}, nil
}

func ReadNDEDNDirs[T any](dirNames []string, regex *regexp.Regexp) ([]FileContent[[]T], error) {
	return readDirs(dirNames, regex, ReadNDEDNFile[T])
}

func WriteNDEDNFile[T any](filePath string, content []T) error {
	buffer := bytes.NewBuffer([]byte{})
	for _, c := range content {
		data, err := edn.Marshal(c)
		if err != nil {
			return err
		}
		_, _ = buffer.Write(data)
		_, _ = buffer.Write([]byte("\n"))
	}
	return WriteFile(filePath, buffer.Bytes())
}
