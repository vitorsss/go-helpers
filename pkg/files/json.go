package files

import (
	"encoding/json"
	"os"
	"regexp"
)

func ReadJSONFile[T any](filePath string) (*FileContent[T], error) {
	fileInfo, err := ReadFileInfo(filePath)
	if err != nil {
		return nil, err
	}

	rawContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var content T
	err = json.Unmarshal(rawContent, &content)
	if err != nil {
		return nil, err
	}

	return &FileContent[T]{
		FileInfo: *fileInfo,
		Content:  content,
	}, nil
}

func ReadJSONDirs[T any](dirNames []string, regex *regexp.Regexp) ([]FileContent[T], error) {
	return readDirs[FileContent[T]](dirNames, regex, ReadJSONFile[T])
}

func WriteJSONFile(filePath string, content interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return WriteFile(filePath, data)
}
