package test

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetAPack(t *testing.T) {
	groupPath := "/group"
	lockPath := "/lock"
	GetPackInSingleGroup(groupPath, lockPath)
	if totalCCUeachIns*totalInstance >= totalPack {
		assert.Equal(t, totalPack, GetNumPackOpened(conns[0], groupPath))
	} else {
		assert.Equal(t, totalCCUeachIns*totalInstance, GetNumPackOpened(conns[0], groupPath))
	}
	assert.Equal(t, DeleteGroup(conns[0], groupPath), nil)
}

func TestGetMultiple(t *testing.T) {
	listGroupPath := GetPackInMultipleGroup()
	for i := 0; i < len(listGroupPath); i++ {
		if totalCCUeachIns*totalInstance >= totalPack {
			assert.Equal(t, totalPack, GetNumPackOpened(conns[0], listGroupPath[i]))
		} else {
			assert.Equal(t, totalCCUeachIns*totalInstance, GetNumPackOpened(conns[0], listGroupPath[i]))
		}
		assert.Equal(t, DeleteGroup(conns[0], listGroupPath[i]), nil)
	}
}
