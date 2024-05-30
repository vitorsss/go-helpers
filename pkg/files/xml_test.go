package files

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
)

func Test_ReadXMLDirs(t *testing.T) {
	fileContents, err := ReadXMLDirs[someType](
		[]string{
			"./testdata",
		},
		regexp.MustCompile(".*\\.xml"),
	)
	if !assertutil.Error(t, nil, err) {
		return
	}

	if !assert.Equal(t, 1, len(fileContents)) {
		return
	}

	assert.Equal(t, someType{
		ID:    1,
		Value: "ttt",
	}, fileContents[0].Content)
}

func Test_WriteXMLFile(t *testing.T) {
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

	filePath := path.Join(dirPath, "somefile.xml")

	err = WriteXMLFile(
		filePath,
		someType{
			ID:    1,
			Value: "bb",
		},
	)
	if !assertutil.Error(t, nil, err) {
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, `<someType><id>1</id><value>bb</value></someType>`, string(data))
}

func Test_ReadNDXMLDirs(t *testing.T) {
	fileContents, err := ReadNDXMLDirs[someType](
		[]string{
			"./testdata",
		},
		regexp.MustCompile(".*\\.ndxml"),
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

func Test_WriteNDXMLFile(t *testing.T) {
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

	filePath := path.Join(dirPath, "somefile.ndxml")

	err = WriteNDXMLFile(
		filePath,
		[]someType{{
			ID:    1,
			Value: "bb",
		}, {
			ID:    3,
			Value: "ss\ndd",
		}},
	)
	if !assertutil.Error(t, nil, err) {
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, `<someType><id>1</id><value>bb</value></someType>
<someType><id>3</id><value>ss&#xA;dd</value></someType>
`, string(data))
}
