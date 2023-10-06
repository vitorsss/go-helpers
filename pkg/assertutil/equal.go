package assertutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func JSONEqual(t *testing.T, a interface{}, b interface{}) bool {
	aData, aErr := json.Marshal(a)
	bData, bErr := json.Marshal(b)
	return Error(t, aErr, bErr) && assert.Equal(t, string(aData), string(bData))
}
