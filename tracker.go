package ratetracker

import (
	"sync"
	"time"
)

//Bucket is simply a bucket
type Bucket struct {
	l        sync.RWMutex
	entries  []time.Time
	tick     time.Time
	duration time.Duration
}

//Len returns the number of entries specified in the
func (bucket *Bucket) Len() int {
	bucket.l.Lock()
	defer bucket.l.Unlock()
	checkAndDeleteOld(bucket)
	//if last tick was earlier than 15 mins ago then set the tick to current time
	return len(bucket.entries)
}
func checkAndDeleteOld(b *Bucket) {
	//this function assumes the bucket is locked elsewhere. If the bucket is not locked before being passed in this function
	//it will likely cause a race condition
	if b.tick.Unix() <= time.Now().Add(-1*b.duration).Unix() {
		b.tick = time.Now()
	}

	offset := -1
	for i, v := range b.entries {
		if v.Unix() <= b.tick.Unix() {
			offset = i
		} else {
			break
		}
	}
	b.entries = b.entries[offset:]
}

//NewBucket creates a new bucket with a timeframe given in seconds
func NewBucket(duration time.Duration) *Bucket {
	return &Bucket{
		entries:  make([]time.Time, 0, 100),
		tick:     time.Now(),
		duration: duration,
	}
}
