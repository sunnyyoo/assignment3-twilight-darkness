package cache

// An FIFO is a fixed-size in-memory cache with first-in first-out eviction
type FIFO struct {
	// whatever fields you want here
	stats *Stats
	len int
	remaining int
	limit int
	kvmap map[string][]byte
	queue []string
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewFifo(limit int) *FIFO {
	fifo := new(FIFO)
	fifo.stats = new(Stats)
	fifo.len = 0
	fifo.limit = limit
	fifo.remaining = limit
	fifo.kvmap = make(map[string][]byte)
	fifo.queue = make([]string, 0)
	return fifo
}

// MaxStorage returns the maximum number of bytes this FIFO can store
func (fifo *FIFO) MaxStorage() int {
	return fifo.limit
}

// RemainingStorage returns the number of unused bytes available in this FIFO
func (fifo *FIFO) RemainingStorage() int {
	return fifo.remaining
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (fifo *FIFO) Get(key string) (value []byte, ok bool) {
	val, ok := fifo.kvmap[key]
	if ok {
		fifo.stats.Hits++
	} else {
		fifo.stats.Misses++
	}
	return val, ok
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {
	val, ok := fifo.kvmap[key]
	if ok {
		fifo.stats.Hits++
	} else {
		fifo.stats.Misses++
	}
	fifo.remaining += len(key) + len(val)
	delete(fifo.kvmap, key)
	return val, ok
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(key string, value []byte) bool {
	if len(key) + len(value) > fifo.limit {
		return false
	}
	val, ok := fifo.kvmap[key]
	if ok {
		fifo.remaining -= len(value) - len(val)
	} else {
		fifo.remaining -= len(value) + len(key)
	}
	fifo.kvmap[key] = value
	for fifo.remaining < 0 {
		fi := fifo.queue[0]
		_, ok = fifo.kvmap[fi]
		if ok {
			fifo.Remove(fi)
		}
	}
	return true
}

// Len returns the number of bindings in the FIFO.
func (fifo *FIFO) Len() int {
	return len(fifo.kvmap)
}

// Stats returns statistics about how many search hits and misses have occurred.
func (fifo *FIFO) Stats() *Stats {
	return fifo.stats
}
