package ratetracker

import (
	"fmt"
	"sync"
	"time"
)

//Bucket is simply a bucket
type Bucket struct {
	l        sync.Mutex
	entries  []time.Time
	duration time.Duration
}

//Len returns the number of entries specified in the
func (b *Bucket) Len() int {
	b.l.Lock()
	defer b.l.Unlock()
	checkAndDeleteOld(b)
	return len(b.entries)
}

//Add adds new entry to the list
func (b *Bucket) Add() {
	b.l.Lock()
	defer b.l.Unlock()
	b.entries = append(b.entries, time.Now())
	checkAndDeleteOld(b)
}

//NewBucket creates a new bucket with a timeframe given in seconds
func NewBucket(duration time.Duration) *Bucket {
	return &Bucket{
		entries:  make([]time.Time, 0, 100),
		duration: duration,
	}
}

//BucketMap type
type BucketMap struct {
	l sync.Mutex
	m map[string]*Bucket
	t time.Duration
}

//NewBucketMap creates a BucketMap
func NewBucketMap(duration time.Duration) *BucketMap {

	return &BucketMap{
		m: make(map[string]*Bucket),
		t: duration,
	}
}

//Entry creates a new entry or updates the existing one
func (m *BucketMap) Entry(key string) {
	m.l.Lock()
	defer m.l.Unlock()
	bucket, ok := m.m[key]
	if ok {
		bucket.Add()
	} else {
		buck := NewBucket(m.t)
		buck.Add()
		m.m[key] = buck
	}
}

//Len returns the number of the entries in specified time for the key
func (m *BucketMap) Len(key string) int {
	m.l.Lock()
	defer m.l.Unlock()
	bucket, ok := m.m[key]
	if ok {
		return bucket.Len()
	}
	return 0

}

func checkAndDeleteOld(b *Bucket) {
	//this function assumes the bucket is locked elsewhere. If the bucket is not locked before being passed in this function
	//it will likely cause a race condition
	if len(b.entries) == 0 {
		return
	}

	offset := -1
	for i, v := range b.entries {
		if v.Unix() < time.Now().Add(-1*b.duration).Unix() {
			offset = i
			fmt.Println("old data")
		} else {
			break
		}
	}
	b.entries = b.entries[offset+1:]
}
