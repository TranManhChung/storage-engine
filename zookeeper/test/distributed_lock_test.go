package test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/zk"
)

func TestGetPack(t *testing.T) {
	var wgIns sync.WaitGroup
	for i := 0; i < totalInstance; i++ {
		wgIns.Add(1)
		go func(i int) { // start totalInstance instance
			defer wgIns.Done()

			var wg sync.WaitGroup
			for j := 0; j < totalCCUeachIns; j++ {
				wg.Add(1)
				go func() { // simulator concurrent user each instance
					defer wg.Done()
					GetPack(conns[i])
				}()
			}
			wg.Wait()
		}(i)
	}
	wgIns.Wait()

	if totalCCUeachIns*totalInstance >= totalPack {
		assert.Equal(t, totalPack, GetNumPackOpened(conns[0]))
	} else {
		assert.Equal(t, totalCCUeachIns*totalInstance, GetNumPackOpened(conns[0]))
	}

	assert.Equal(t, DeleteGroup(conns[0]), nil)
}

func GetPack(c *zk.Conn) {
	var data []byte
	var stt *zk.Stat
	var err error
	group := Group{Packages: make([]int, totalPack)}

	l := zk.NewLock(c, lockPath, acl)
	l.Lock()
	defer l.Unlock()

	if data, stt, err = c.Get(groupPath); err != nil {
		if data, err = json.Marshal(group); err != nil {
			fmt.Println("MARSHAL DATA FAILED !!!")
		}
		if _, err = c.Create(groupPath, data, 0, acl); err != nil {
			fmt.Println("ERROR : ", err)
		}
	}

	if err = json.Unmarshal(data, &group); err != nil {
		fmt.Println("ERROR : ", err)
	}
	i := 0
	for ; i < len(group.Packages); i++ {
		if group.Packages[i] != -1 {
			group.Packages[i] = -1
			break
		}
	}
	if i >= len(group.Packages) {
		return
	}
	if data, err = json.Marshal(group); err != nil {
		fmt.Println("ERROR : ", err)
		return
	}

	if _, err = c.Set(groupPath, data, stt.Version); err != nil {
		fmt.Println("ERROR : ", err)
	}
}

func GetNumPackOpened(c *zk.Conn) int {
	group := Group{Packages: make([]int, totalPack)}
	data, _, _ := c.Get(groupPath)
	json.Unmarshal(data, &group)
	result := 0
	for i := 0; i < len(group.Packages); i++ {
		if group.Packages[i] == -1 {
			result++
		}
	}
	return result
}

func DeleteGroup(c *zk.Conn) error {
	_, stt, err := c.Get(groupPath)
	if err != nil {
		return err
	}
	err = c.Delete(groupPath, stt.Version)
	if err != nil {
		return err
	}
	return nil
}
