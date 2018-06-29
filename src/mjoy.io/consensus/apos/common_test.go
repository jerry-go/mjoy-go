package apos

import (
	"math/big"
	"testing"
	"container/heap"
	"fmt"
)

type testObj struct {
	aa    int
	bb    *big.Int
}
func TestPg(t *testing.T) {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)

	for i:=0; i<10; i++ {
		to0 := &testObj{i, big.NewInt(int64(i))}
		item0 := &pqItem{to0, to0.bb}
		heap.Push(&pq, item0)
	}

	for i:=0; i<10; i++ {
		to1 := &testObj{100 -i, big.NewInt(int64(100 -i +1))}
		item1 := &pqItem{to1, to1.bb}
		heap.Push(&pq, item1)
	}
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*pqItem)
		fmt.Println("output", item.priority, item.value.(*testObj).aa, item.value.(*testObj).bb)
	}
}

func TestPg1(t *testing.T) {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)

	for i:=0; i<10; i++ {
		to0 := &testObj{i, big.NewInt(int64(i))}
		item0 := &pqItem{to0, to0.bb}
		heap.Push(&pq, item0)
	}

	for i:=0; i<10; i++ {
		to1 := &testObj{100 -i, big.NewInt(int64(100 -i +1))}
		item1 := &pqItem{to1, to1.bb}
		heap.Push(&pq, item1)
	}
	for i:= 0; i < pq.Len(); i++ {
		fmt.Println("output", pq[i].priority, pq[i].value.(*testObj).aa, pq[i].value.(*testObj).bb)
	}
}
