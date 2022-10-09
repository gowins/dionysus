package algs

import (
	"testing"
)

func TestPop(t *testing.T) {
	pq := NewPQ()
	if _, ok := pq.Pop(); ok {
		t.Error("should return nil & false")
	}
}

func TestGetItem(t *testing.T) {
	pq := NewPQ()
	itemsIn := map[string]int{
		"test2": 2,
		"test1": 1,
		"test3": 3,
		"test0": 0,
	}
	itemsValueOut := []string{
		"test0",
		"test1",
		"test2",
		"test3",
	}
	for key, value := range itemsIn {
		item := NewItem(key, value)
		pq.Push(item)
	}
	for i := 0; i < len(itemsValueOut); i++ {
		item, ok := pq.Pop()
		if !ok || item.Value() != itemsValueOut[i] {
			t.Errorf("expect ok: true, get: %v,expect value: %v, get value: %v", ok, itemsValueOut[i], item.Value())

		}
	}
}

func TestUpdate(t *testing.T) {
	pq := NewPQ()
	itemsIn := map[string]int{
		"test2": 2,
		"test1": 1,
		"test3": 3,
		"test0": 0,
	}
	itemsValueOut := []string{
		"test1",
		"test2",
		"test3",
		"test5",
	}

	var item *Item
	for key, value := range itemsIn {
		it := NewItem(key, value)
		if key == "test0" {
			item = it
		}
		pq.Push(it)
	}

	// update the item{"test0": 0} to item{"test5": 5}
	pq.Update(item, "test5", 5)

	for i := 0; i < len(itemsValueOut); i++ {
		item, ok := pq.Pop()
		if !ok || item.Value() != itemsValueOut[i] {
			t.Errorf("expect ok: true, get: %v,expect value: %v, get value: %v", ok, itemsValueOut[i], item.Value())
		}
	}
}
