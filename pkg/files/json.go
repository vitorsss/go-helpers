package files

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"regexp"

	"github.com/pkg/errors"
)

func ReadJSONFile[T any](filePath string) (*FileContent[T], error) {
	fileInfo, err := ReadFileInfo(filePath)
	if err != nil {
		return nil, err
	}

	rawContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	var content T
	err = json.Unmarshal(rawContent, &content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal content")
	}

	return &FileContent[T]{
		FileInfo: *fileInfo,
		Content:  content,
	}, nil
}

func ReadJSONDirs[T any](dirNames []string, regex *regexp.Regexp) ([]FileContent[T], error) {
	return readDirs(dirNames, regex, ReadJSONFile[T])
}

func WriteJSONFile(filePath string, content interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return errors.Wrap(err, "failed to marshal content")
	}
	return WriteFile(filePath, data)
}

func ReadNDJSONFile[T any](filePath string) (*FileContent[[]T], error) {
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
		err = json.Unmarshal([]byte(scanner.Text()), &contentLine)
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

func ReadNDJSONDirs[T any](dirNames []string, regex *regexp.Regexp) ([]FileContent[[]T], error) {
	return readDirs(dirNames, regex, ReadNDJSONFile[T])
}

func WriteNDJSONFile[T any](filePath string, content []T) error {
	buffer := bytes.NewBuffer([]byte{})
	for _, c := range content {
		data, err := json.Marshal(c)
		if err != nil {
			return err
		}
		_, _ = buffer.Write(data)
		_, _ = buffer.Write([]byte("\n"))
	}
	return WriteFile(filePath, buffer.Bytes())
}
