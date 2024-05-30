package files

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
)

type someType struct {
	ID    int    `json:"id" edn:"id" xml:"id"`
	Value string `json:"value" edn:"tt/value" xml:"value"`
}

func Test_CreateDatedOutput(t *testing.T) {
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

	w, err := CreateDatedOutput(dirPath, "somefilename", "ext")
	if !assertutil.Error(t, nil, err) {
		return
	}
	now := time.Now()
	_, err = w.Write([]byte("test"))
	if err != nil {
		panic(err)
	}

	err = w.Close()
	if err != nil {
		panic(err)
	}

	dirEntries, err := os.ReadDir(path.Join(dirPath, now.Format(outputTimeFormat)))
	if err != nil {
		panic(err)
	}

	if !assert.Equal(t, 1, len(dirEntries)) {
		return
	}

	dirEntry := dirEntries[0]
	assert.True(t, strings.HasPrefix(dirEntry.Name(), "somefilename"))
}
