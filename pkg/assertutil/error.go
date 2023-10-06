package assertutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Error(t *testing.T, wantErr error, err error) bool {
	if wantErr == nil {
		return assert.Nil(t, err)
	} else if assert.NotNil(t, err) {
		return assert.Equal(t, wantErr.Error(), err.Error())
	}
	return false
}
