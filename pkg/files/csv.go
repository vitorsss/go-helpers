package files

import (
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

func ReadCSVFile(filePath string, comma rune, hasHeaders bool) (*FileContent[[]map[string]string], error) {
	fileInfo, err := ReadFileInfo(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = comma

	content := []map[string]string{}

	headerByIndex := map[int]string{}
	if hasHeaders {
		headers, err := csvReader.Read()
		if err != nil {
			return nil, errors.Wrap(err, "failed to read line")
		}
		for idx, header := range headers {
			headerByIndex[idx] = header
		}
	}

	for line, err := csvReader.Read(); err == nil; line, err = csvReader.Read() {
		lineContent := map[string]string{}
		for idx, cell := range line {
			var name string
			var ok bool
			if name, ok = headerByIndex[idx]; !ok {
				name = strconv.Itoa(idx)
			}
			lineContent[name] = cell
		}
		content = append(content, lineContent)
	}
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "failed to read line")
	}

	return &FileContent[[]map[string]string]{
		FileInfo: *fileInfo,
		Content:  content,
	}, nil
}

func ReadCSVDirs(dirNames []string, comma rune, hasHeaders bool, regex *regexp.Regexp) ([]FileContent[[]map[string]string], error) {
	return readDirs(dirNames, regex, func(filePath string) (*FileContent[[]map[string]string], error) {
		return ReadCSVFile(filePath, comma, hasHeaders)
	})
}
