package cache

type LinkedListNode struct {
	next *LinkedListNode
	prev *LinkedListNode
	data string
}
// An LRU is a fixed-size in-memory cache with least-recently-used eviction
type LRU struct {
	// whatever fields you want here
	stats *Stats
	len int
	remaining int
	limit int
	kvmap map[string][]byte
	knmap map[string]*LinkedListNode
	head *LinkedListNode
	tail *LinkedListNode
}

func (lru *LRU) RemoveKey(key string) {
	node, ok := lru.knmap[key]
	if ok {
		delete(lru.knmap, key)
		if node.prev != nil {
			node.prev.next = node.next
		} else {
			lru.head = node.next
		}
		if node.next != nil {
			node.next.prev = node.prev
		} else {
			lru.tail = node.prev
		}
	}
}

func (lru *LRU) AddKey(key string) {
	newNode := new(LinkedListNode)
	lru.knmap[key] = newNode
	newNode.data = key
	if lru.tail != nil {
		newNode.prev = lru.tail
		lru.tail.next = newNode
		lru.tail = newNode
	} else {
		lru.head = newNode
		lru.tail = newNode
	}
}

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {
	lru := new(LRU)
	lru.stats = new(Stats)
	lru.len = 0
	lru.limit = limit
	lru.remaining = limit
	lru.kvmap = make(map[string][]byte)
	lru.knmap = make(map[string]*LinkedListNode)
	return lru
}

// MaxStorage returns the maximum number of bytes this LRU can store
func (lru *LRU) MaxStorage() int {
	return lru.limit
}

// RemainingStorage returns the number of unused bytes available in this LRU
func (lru *LRU) RemainingStorage() int {
	return lru.remaining
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value []byte, ok bool) {
	val, ok := lru.kvmap[key]
	if ok {
		lru.stats.Hits++
		lru.RemoveKey(key)
		lru.AddKey(key)
	} else {
		lru.stats.Misses++
	}
	return val, ok
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value []byte, ok bool) {
	val, ok := lru.kvmap[key]
	if ok {
		lru.remaining += len(key) + len(val)
	}
	lru.RemoveKey(key)
	return val, ok
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(key string, value []byte) bool {
	if len(key) + len(value) > lru.limit {
		return false
	}
	val, ok := lru.kvmap[key]
	if ok {
		lru.remaining -= len(value) - len(val)
	} else {
		lru.remaining -= len(value) + len(key)
	}
	lru.kvmap[key] = value
	lru.AddKey(key)
	for lru.remaining < 0 {
		lr := lru.head
		lru.Remove(lr.data)
	}
	return true
}

// Len returns the number of bindings in the LRU.
func (lru *LRU) Len() int {
	return len(lru.kvmap)
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lru *LRU) Stats() *Stats {
	return lru.stats
}
