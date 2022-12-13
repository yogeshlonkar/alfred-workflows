package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xujiajun/nutsdb"
)

const (
	minusOneDay = time.Hour * 24 * -1
	segmentSize = 16 * nutsdb.MB
)

var (
	lastSyncKey = []byte("last_sync")
	yesterday   = time.Now().Add(minusOneDay)
)

type cache[K any] struct {
	*nutsdb.DB
	bucket, path string
}

func (c *cache[K]) Close() error {
	err := c.DB.Close()
	if err != nil {
		return fmt.Errorf("error closing cache: %w", err)
	}
	return nil
}

func (c *cache[K]) Delete() (err error) {
	if c.DB != nil {
		if err = c.Close(); err != nil {
			return err
		}
	}
	if err = os.RemoveAll(c.path); err != nil {
		return fmt.Errorf("error removing cache %s: %w", c.path, err)
	}
	return nil
}

func (c *cache[K]) Exists() bool {
	return c.KeyCount > 0
}

func (c *cache[K]) Get(key string) (v K, err error) {
	err = c.View(
		func(tx *nutsdb.Tx) error {
			var e *nutsdb.Entry
			e, err = tx.Get(c.bucket, []byte(key))
			if err != nil {
				return fmt.Errorf("error getting entity on key %s: %w", key, err)
			}
			if err = json.Unmarshal(e.Value, &v); err != nil {
				return fmt.Errorf("error decoding JSON: %w", err)
			}
			return nil
		})
	if errors.Is(err, nutsdb.ErrBucketNotFound) {
		return v, nil
	}
	return
}

func (c *cache[K]) GetJSON(key string) (json []byte, err error) {
	err = c.View(
		func(tx *nutsdb.Tx) error {
			var e *nutsdb.Entry
			e, err = tx.Get(c.bucket, []byte(key))
			if err != nil {
				return fmt.Errorf("error getting JSON entity on key %s: %w", key, err)
			}
			json = e.Value
			return nil
		})
	if errors.Is(err, nutsdb.ErrBucketNotFound) {
		return nil, nil
	}
	return
}

func (c *cache[K]) Open() (err error) {
	c.DB, err = nutsdb.Open(
		nutsdb.Options{
			EntryIdxMode: nutsdb.HintKeyValAndRAMIdxMode,
			SegmentSize:  segmentSize,
			NodeNum:      1,
			RWMode:       nutsdb.FileIO,
			SyncEnable:   true,
		},
		nutsdb.WithDir(c.path),
	)
	if err != nil {
		return fmt.Errorf("error opening cache: %w", err)
	}
	return nil
}

func (c *cache[K]) Path() string {
	return c.path
}

func (c *cache[K]) Save(key string, v K) (data []byte, err error) {
	if data, err = json.Marshal(v); err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}
	err = c.Update(
		func(tx *nutsdb.Tx) error {
			err = tx.Put(c.bucket, lastSyncKey, []byte(time.Now().Format(time.RFC3339)), 0)
			if err != nil {
				return fmt.Errorf("error saving last_sync: %w", err)
			}
			err = tx.Put(c.bucket, []byte(key), data, 0)
			if err != nil {
				return fmt.Errorf("error saving entity on key %s: %w", key, err)
			}
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("error sving entity on key %s: %w", key, err)
	}
	return
}

func (c *cache[K]) Stale() bool {
	var lastSync time.Time
	err := c.View(
		func(tx *nutsdb.Tx) error {
			e, err := tx.Get(c.bucket, lastSyncKey)
			if err != nil {
				return fmt.Errorf("error getting last_sync: %w", err)
			}
			lastSync, err = time.Parse(time.RFC3339, string(e.Value))
			if err != nil {
				return fmt.Errorf("error parsing last_sync: %w", err)
			}
			return nil
		})
	if err != nil {
		log.Warn().Err(err).Msg("returning stale=true")
		return true
	}
	return lastSync.Before(yesterday)
}
