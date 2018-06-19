package apos

import (
	"testing"
	"time"
	"fmt"
)

func TestAposRunning(t *testing.T){
	an := newAllNodeManager()
	an.init()
	for{
		time.Sleep(1*time.Second)
		fmt.Println("apos_test doing....")
	}
}

