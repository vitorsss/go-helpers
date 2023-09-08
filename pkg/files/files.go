package files

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"time"
)

const (
	outputTimeFormat = "2006/01/02"
)

type FileInfo struct {
	Dir     string
	Name    string
	ModTime time.Time
}

type FileContent[T any] struct {
	FileInfo
	Content T
}

func ReadFileInfo(filePath string) (*FileInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Dir:     path.Dir(filePath),
		Name:    fileInfo.Name(),
		ModTime: fileInfo.ModTime(),
	}, nil
}

func ReadDirsFileInfos(dirNames []string, regex *regexp.Regexp) ([]FileInfo, error) {
	return readDirs[FileInfo](dirNames, regex, ReadFileInfo)
}

func readDirs[T any](dirNames []string, regex *regexp.Regexp, readFileFn func(filePath string) (*T, error)) ([]T, error) {
	if len(dirNames) == 0 {
		return nil, nil
	}
	result := []T{}
	subDirs := []string{}
	for _, dirName := range dirNames {
		dirEntries, err := os.ReadDir(dirName)
		if err != nil {
			return nil, err
		}
		for _, dirEntry := range dirEntries {
			fileName := path.Join(dirName, dirEntry.Name())
			if dirEntry.IsDir() {
				subDirs = append(subDirs, fileName)
			} else if regex.MatchString(fileName) {
				fileContent, err := readFileFn(fileName)
				if err != nil {
					return nil, err
				}
				result = append(result, *fileContent)
			}
		}
	}

	subFilesContent, err := readDirs[T](subDirs, regex, readFileFn)
	if err != nil {
		return nil, err
	}
	result = append(result, subFilesContent...)
	return result, nil
}

func CreateDatedOutput(baseDir string, fileName string, extension string) (io.WriteCloser, error) {
	now := time.Now()
	filePath := path.Join(
		baseDir,
		now.Format(outputTimeFormat),
		fmt.Sprintf("%s-%s.%s",
			fileName,
			now.Format(time.TimeOnly),
			extension,
		),
	)
	err := os.MkdirAll(path.Dir(filePath), 0o777)
	if err != nil {
		return nil, err
	}
	return os.Create(
		filePath,
	)
}

func WriteFile(filePath string, content []byte) error {
	err := os.MkdirAll(path.Dir(filePath), 0o777)
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.WriteFile(filePath, content, 0o444)
}
