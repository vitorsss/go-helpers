package files

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
)

func Test_ReadEDNDirs(t *testing.T) {
	fileContents, err := ReadEDNDirs[[]someType](
		[]string{
			"./testdata",
		},
		regexp.MustCompile(".*\\.edn"),
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

func Test_WriteEDNFile(t *testing.T) {
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

	filePath := path.Join(dirPath, "somefile.edn")

	err = WriteEDNFile(
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

	assert.Equal(t, `[{:id 1 :tt/value"bb"}]`, string(data))
}

func Test_ReadNDEDNDirs(t *testing.T) {
	fileContents, err := ReadNDEDNDirs[someType](
		[]string{
			"./testdata",
		},
		regexp.MustCompile(".*\\.ndedn"),
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

func Test_WriteNDEDNFile(t *testing.T) {
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

	filePath := path.Join(dirPath, "somefile.ndedn")

	err = WriteNDEDNFile(
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

	assert.Equal(t, `{:id 1 :tt/value"bb"}
{:id 3 :tt/value"ss"}
`, string(data))
}
