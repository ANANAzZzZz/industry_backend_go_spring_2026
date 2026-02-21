package main

import (
	"container/list"
)

type LRUCache[K comparable, V any] struct {
	capacity              int
	cachedMap             map[K]V
	linkedListElementsMap map[K]*list.Element
	cachedLinkedList      *list.List
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity:              capacity,
		cachedMap:             make(map[K]V),
		linkedListElementsMap: make(map[K]*list.Element),
		cachedLinkedList:      list.New(),
	}
}

func (L *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	var zeroValue V

	if L.capacity == 0 {
		return zeroValue, false
	}

	foundValue, isExists := L.cachedMap[key]

	if isExists {
		L.makeMRU(key)

		return foundValue, true
	}
	return zeroValue, false
}

func (L *LRUCache[K, V]) Set(key K, value V) {
	if L.capacity == 0 {
		return
	}
	_, isExists := L.cachedMap[key]
	L.cachedMap[key] = value

	if isExists {
		L.makeMRU(key)
		return
	}
	L.processCacheMiss(key)
	L.removeLRU()
}

func (L *LRUCache[K, V]) processCacheMiss(key K) {
	L.cachedLinkedList.PushFront(key)
	L.linkedListElementsMap[key] = L.cachedLinkedList.Front()
}

func (L *LRUCache[K, V]) removeLRU() {
	if tail := L.cachedLinkedList.Back(); L.cachedLinkedList.Len() > L.capacity && tail != nil {
		delete(L.cachedMap, tail.Value.(K))
		delete(L.linkedListElementsMap, tail.Value.(K))

		L.cachedLinkedList.Remove(tail)
	}
}

func (L *LRUCache[K, V]) makeMRU(key K) {
	accessedListElement := L.linkedListElementsMap[key]
	L.cachedLinkedList.MoveToFront(accessedListElement)
}

// Space complexity - O(n)

// Set time complexity - amortized O(1)
// Get time complexity - amortized O(1)
