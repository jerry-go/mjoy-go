package apos

import (
	"testing"
	"sort"
	"fmt"
)

func TestBpWithPriorityHeap(t *testing.T){
	hp := make(BpWithPriorityHeap , 0)

	for i:= 0;i < 20;i++{
		bp := new(BpWithPriority)
		bp.j = i

		hp.Push(bp)
	}

	sort.Sort(hp)
	for index , v := range hp{
		fmt.Println("index :" , index , " j :" , v.j)
	}
}




