package api

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestGetAllDept(t *testing.T) {
	res, err := CallGetAllDept()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, res.GetMapArray("deptList"))
}

func TestCallGetAllMajor(t *testing.T) {
	res, err := CallGetAllMajor()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, res.GetMapArray("majorList"))
}
