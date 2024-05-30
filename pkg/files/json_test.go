package files

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
)

func Test_ReadJSONDirs(t *testing.T) {
	fileContents, err := ReadJSONDirs[[]someType](
		[]string{
			"./testdata",
		},
		regexp.MustCompile(".*\\.json"),
	)
	if !assertutil.Error(t, nil, err) {
		return
	}

	if !assert.Equal(t, 1, len(fileContents)) {
		return
	}

	assert.Equal(t, []someType{
		{
			ID:    1,
			Value: "ttt",
		},
		{
			ID:    2,
			Value: "bbbb",
		},
	}, fileContents[0].Content)
}

func Test_WriteJSONFile(t *testing.T) {
	dirPath, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dirPath)
		if err != nil {
			panic(err)
		}
	}()

	filePath := path.Join(dirPath, "somefile.json")

	err = WriteJSONFile(
		filePath,
		[]someType{{
			ID:    1,
			Value: "bb",
		}},
	)
	if !assertutil.Error(t, nil, err) {
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, `[{"id":1,"value":"bb"}]`, string(data))
}

func Test_ReadNDJSONDirs(t *testing.T) {
	fileContents, err := ReadNDJSONDirs[someType](
		[]string{
			"./testdata",
		},
		regexp.MustCompile(".*\\.ndjson"),
	)
	if !assertutil.Error(t, nil, err) {
		return
	}

	if !assert.Equal(t, 1, len(fileContents)) {
		return
	}

	assert.Equal(t, []someType{
		{
			ID:    1,
			Value: "ii",
		},
		{
			ID:    2,
			Value: "kk",
		},
	}, fileContents[0].Content)
}

func Test_WriteNDJSONFile(t *testing.T) {
	dirPath, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dirPath)
		if err != nil {
			panic(err)
		}
	}()

	filePath := path.Join(dirPath, "somefile.ndjson")

	err = WriteNDJSONFile(
		filePath,
		[]someType{{
			ID:    1,
			Value: "bb",
		}, {
			ID:    3,
			Value: "ss",
		}},
	)
	if !assertutil.Error(t, nil, err) {
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, `{"id":1,"value":"bb"}
{"id":3,"value":"ss"}
`, string(data))
}
