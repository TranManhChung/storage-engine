package test

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/zk"
)

const (
	perms           = 0x1f
	totalPack       = 1000
	totalCCUeachIns = 320
	totalInstance   = 3
	numGroup        = 2
)

var (
	acl   = []zk.ACL{{perms, "world", "anyone"}}
	conns = make([]*zk.Conn, totalInstance) // create a connection for each instance
)

type Group struct {
	Packages []int `json:"package"`
}

func TestMain(m *testing.M) {
	var err error
	// new connection to zk
	for i := 0; i < len(conns); i++ {
		if conns[i], _, err = zk.Connect([]string{"172.17.0.2"}, time.Second); err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}

func GetPackInSingleGroup(groupPath, lockPath string) {
	var wgIns sync.WaitGroup
	for i := 0; i < totalInstance; i++ {
		wgIns.Add(1)
		go func(i int) { // start totalInstance instance
			defer wgIns.Done()

			var wg sync.WaitGroup
			for j := 0; j < totalCCUeachIns; j++ {
				wg.Add(1)
				go func() { // simulator concurrent user each instance
					start := time.Now()
					defer wg.Done()
					GetPack(conns[i], groupPath, lockPath)
					fmt.Println(groupPath, " - ", time.Since(start))
				}()
			}
			wg.Wait()
		}(i)
	}
	wgIns.Wait()
}

func GetPackInMultipleGroup() []string {
	listGroupPath := make([]string, numGroup)
	var wg sync.WaitGroup
	for i := 0; i < numGroup; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			listGroupPath[i] = fmt.Sprintf("/group%d", i)
			GetPackInSingleGroup(listGroupPath[i], fmt.Sprintf("/lock%d", i))
		}(i)
	}
	wg.Wait()
	return listGroupPath
}

func GetPack(c *zk.Conn, groupPath, lockPath string) {
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

func GetNumPackOpened(c *zk.Conn, groupPath string) int {
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

func DeleteGroup(c *zk.Conn, groupPath string) error {
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
