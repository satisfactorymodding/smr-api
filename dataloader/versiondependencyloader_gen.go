// Code generated by github.com/vektah/dataloaden, DO NOT EDIT.

package dataloader

import (
	"sync"
	"time"

	"github.com/satisfactorymodding/smr-api/db/postgres"
)

// VersionDependencyLoaderConfig captures the config to create a new VersionDependencyLoader
type VersionDependencyLoaderConfig struct {
	// Fetch is a method that provides the data for the loader
	Fetch func(keys []string) ([][]postgres.VersionDependency, []error)

	// Wait is how long wait before sending a batch
	Wait time.Duration

	// MaxBatch will limit the maximum number of keys to send in one batch, 0 = not limit
	MaxBatch int
}

// NewVersionDependencyLoader creates a new VersionDependencyLoader given a fetch, wait, and maxBatch
func NewVersionDependencyLoader(config VersionDependencyLoaderConfig) *VersionDependencyLoader {
	return &VersionDependencyLoader{
		fetch:    config.Fetch,
		wait:     config.Wait,
		maxBatch: config.MaxBatch,
	}
}

// VersionDependencyLoader batches and caches requests
type VersionDependencyLoader struct {
	fetch    func(keys []string) ([][]postgres.VersionDependency, []error)
	cache    map[string][]postgres.VersionDependency
	batch    *versionDependencyLoaderBatch
	wait     time.Duration
	maxBatch int
	mu       sync.Mutex
}

type versionDependencyLoaderBatch struct {
	done    chan struct{}
	keys    []string
	data    [][]postgres.VersionDependency
	error   []error
	closing bool
}

// Load a VersionDependency by key, batching and caching will be applied automatically
func (l *VersionDependencyLoader) Load(key string) ([]postgres.VersionDependency, error) {
	return l.LoadThunk(key)()
}

// LoadThunk returns a function that when called will block waiting for a VersionDependency.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *VersionDependencyLoader) LoadThunk(key string) func() ([]postgres.VersionDependency, error) {
	l.mu.Lock()
	if it, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]postgres.VersionDependency, error) {
			return it, nil
		}
	}
	if l.batch == nil {
		l.batch = &versionDependencyLoaderBatch{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()

	return func() ([]postgres.VersionDependency, error) {
		<-batch.done

		var data []postgres.VersionDependency
		if pos < len(batch.data) {
			data = batch.data[pos]
		}

		var err error
		// its convenient to be able to return a single error for everything
		if len(batch.error) == 1 {
			err = batch.error[0]
		} else if batch.error != nil {
			err = batch.error[pos]
		}

		if err == nil {
			l.mu.Lock()
			l.unsafeSet(key, data)
			l.mu.Unlock()
		}

		return data, err
	}
}

// LoadAll fetches many keys at once. It will be broken into appropriate sized
// sub batches depending on how the loader is configured
func (l *VersionDependencyLoader) LoadAll(keys []string) ([][]postgres.VersionDependency, []error) {
	results := make([]func() ([]postgres.VersionDependency, error), len(keys))

	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}

	versionDependencys := make([][]postgres.VersionDependency, len(keys))
	errors := make([]error, len(keys))
	for i, thunk := range results {
		versionDependencys[i], errors[i] = thunk()
	}
	return versionDependencys, errors
}

// LoadAllThunk returns a function that when called will block waiting for a VersionDependencys.
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *VersionDependencyLoader) LoadAllThunk(keys []string) func() ([][]postgres.VersionDependency, []error) {
	results := make([]func() ([]postgres.VersionDependency, error), len(keys))
	for i, key := range keys {
		results[i] = l.LoadThunk(key)
	}
	return func() ([][]postgres.VersionDependency, []error) {
		versionDependencys := make([][]postgres.VersionDependency, len(keys))
		errors := make([]error, len(keys))
		for i, thunk := range results {
			versionDependencys[i], errors[i] = thunk()
		}
		return versionDependencys, errors
	}
}

// Prime the cache with the provided key and value. If the key already exists, no change is made
// and false is returned.
// (To forcefully prime the cache, clear the key first with loader.clear(key).prime(key, value).)
func (l *VersionDependencyLoader) Prime(key string, value []postgres.VersionDependency) bool {
	l.mu.Lock()
	var found bool
	if _, found = l.cache[key]; !found {
		// make a copy when writing to the cache, its easy to pass a pointer in from a loop var
		// and end up with the whole cache pointing to the same value.
		cpy := make([]postgres.VersionDependency, len(value))
		copy(cpy, value)
		l.unsafeSet(key, cpy)
	}
	l.mu.Unlock()
	return !found
}

// Clear the value at key from the cache, if it exists
func (l *VersionDependencyLoader) Clear(key string) {
	l.mu.Lock()
	delete(l.cache, key)
	l.mu.Unlock()
}

func (l *VersionDependencyLoader) unsafeSet(key string, value []postgres.VersionDependency) {
	if l.cache == nil {
		l.cache = map[string][]postgres.VersionDependency{}
	}
	l.cache[key] = value
}

// keyIndex will return the location of the key in the batch, if its not found
// it will add the key to the batch
func (b *versionDependencyLoaderBatch) keyIndex(l *VersionDependencyLoader, key string) int {
	for i, existingKey := range b.keys {
		if key == existingKey {
			return i
		}
	}

	pos := len(b.keys)
	b.keys = append(b.keys, key)
	if pos == 0 {
		go b.startTimer(l)
	}

	if l.maxBatch != 0 && pos >= l.maxBatch-1 {
		if !b.closing {
			b.closing = true
			l.batch = nil
			go b.end(l)
		}
	}

	return pos
}

func (b *versionDependencyLoaderBatch) startTimer(l *VersionDependencyLoader) {
	time.Sleep(l.wait)
	l.mu.Lock()

	// we must have hit a batch limit and are already finalizing this batch
	if b.closing {
		l.mu.Unlock()
		return
	}

	l.batch = nil
	l.mu.Unlock()

	b.end(l)
}

func (b *versionDependencyLoaderBatch) end(l *VersionDependencyLoader) {
	b.data, b.error = l.fetch(b.keys)
	close(b.done)
}
