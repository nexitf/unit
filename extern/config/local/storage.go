package local

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	ErrDuplicateKeyName = errors.New("duplicate key name")
)

type Local struct {
	mutex     sync.Mutex
	storage   map[string]func(value string)
	dir       string
	ext       string
	watcher   *fsnotify.Watcher
	cancelCtx context.CancelFunc
}

// NewLocalStorage
func NewLocalStorage(dir, ext string) (l *Local, err error) {
	l = &Local{
		dir:     dir,
		ext:     ext,
		storage: make(map[string]func(value string)),
	}
	// Trim suffix
	if strings.HasSuffix(l.dir, string(os.PathSeparator)) {
		l.dir = strings.TrimSuffix(l.dir, string(os.PathSeparator))
	}
	// Fill a dot
	if !strings.HasPrefix(l.ext, ".") {
		l.ext = "." + l.ext
	}
	// Create a file watcher
	l.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = l.watcher.Add(dir)
	if err != nil {
		return nil, err
	}

	var (
		ctx context.Context
		wg  = new(sync.WaitGroup)
	)
	ctx, l.cancelCtx = context.WithCancel(context.Background())

	wg.Add(1)
	go l.background(ctx, wg)

	wg.Wait()

	return
}

// Name implements config.Storage.
func (l *Local) Name() string {
	return "Local"
}

// Version implements config.Storage.
func (l *Local) Version() string {
	return "1.0.0"
}

// Package implements config.Storage.
func (l *Local) Package() string {
	return "github.com/nexitf/unit/extern/config/local"
}

// Watch implements config.Storage.
func (l *Local) Watch(name string, update func(value string)) (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if _, ok := l.storage[name]; ok {
		return ErrDuplicateKeyName
	}
	l.storage[name] = update
	value, err := os.ReadFile(l.dir + "/" + name + l.ext)
	if err != nil {
		update("")
	} else {
		update(string(value))
	}
	return nil
}

// Close
func (l *Local) Close() {
	if l.cancelCtx != nil {
		l.cancelCtx()
	}
}

// key
func (l *Local) key(file string) string {
	filename := filepath.Base(file)
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

// update
func (l *Local) update(key, file string) (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	fn, ok := l.storage[key]
	if !ok {
		return nil
	}
	if path.Ext(file) != l.ext {
		return nil
	}
	// File changed
	buff, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	fn(string(buff))
	return
}

// delete
func (l *Local) delete(key, file string) (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	fn, ok := l.storage[key]
	if !ok {
		return nil
	}
	if path.Ext(file) != l.ext {
		return nil
	}
	fn("")
	return
}

// background
func (l *Local) background(ctx context.Context, wg *sync.WaitGroup) {
	wg.Done()

	for {
		select {
		// Cancel context
		case <-ctx.Done():
			return
		// Watch events
		case event, ok := <-l.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) {
				l.update(l.key(event.Name), event.Name)
			}
			if event.Has(fsnotify.Rename) {
				l.delete(l.key(event.Name), event.Name)
			}
			if event.Has(fsnotify.Write) {
				l.update(l.key(event.Name), event.Name)
			}
			if event.Has(fsnotify.Remove) {
				l.delete(l.key(event.Name), event.Name)
			}
		// Throw error
		case err, ok := <-l.watcher.Errors:
			if !ok {
				return
			}
			_ = err
		}
	}
}
