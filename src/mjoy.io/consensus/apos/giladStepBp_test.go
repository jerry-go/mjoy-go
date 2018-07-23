package apos

import (
	"testing"
	"fmt"
	"container/heap"
)

func TestBpWithPriorityHeap(t *testing.T){
	hp := make(BpWithPriorityHeap , 0)
	heap.Init(&hp)
	for i:= 0;i < 20;i++{
		bp := new(BpWithPriority)
		bp.j = i
		heap.Push(&hp,bp)

	}
	//we just need the first one
	//sort.Sort(hp)
	for index , v := range hp{
		fmt.Println("index :" , index , " j :" , v.j)
	}
}





