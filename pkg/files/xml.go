package files

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"os"
	"regexp"

	"github.com/pkg/errors"
)

func ReadXMLFile[T any](filePath string) (*FileContent[T], error) {
	fileInfo, err := ReadFileInfo(filePath)
	if err != nil {
		return nil, err
	}

	rawContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	var content T
	err = xml.Unmarshal(rawContent, &content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal content")
	}

	return &FileContent[T]{
		FileInfo: *fileInfo,
		Content:  content,
	}, nil
}

func ReadXMLDirs[T any](dirNames []string, regex *regexp.Regexp) ([]FileContent[T], error) {
	return readDirs(dirNames, regex, ReadXMLFile[T])
}

func WriteXMLFile(filePath string, content interface{}) error {
	data, err := xml.Marshal(content)
	if err != nil {
		return errors.Wrap(err, "failed to marshal content")
	}
	return WriteFile(filePath, data)
}

func ReadNDXMLFile[T any](filePath string) (*FileContent[[]T], error) {
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
		err = xml.Unmarshal([]byte(scanner.Text()), &contentLine)
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

func ReadNDXMLDirs[T any](dirNames []string, regex *regexp.Regexp) ([]FileContent[[]T], error) {
	return readDirs(dirNames, regex, ReadNDXMLFile[T])
}

func WriteNDXMLFile[T any](filePath string, content []T) error {
	buffer := bytes.NewBuffer([]byte{})
	for _, c := range content {
		data, err := xml.Marshal(c)
		if err != nil {
			return err
		}
		_, _ = buffer.Write(data)
		_, _ = buffer.Write([]byte("\n"))
	}
	return WriteFile(filePath, buffer.Bytes())
}
