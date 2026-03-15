package editor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"sort"
	"sync"
)

type mediaInventoryEntry struct {
	hash     string
	partPath string
}

func (e *PresentationEditor) verifyMediaInventoryChecksumsParallel() error {
	entries := e.snapshotMediaInventoryEntries()
	if len(entries) == 0 {
		return nil
	}

	workerCount := runtime.GOMAXPROCS(0)
	workerCount = min(max(workerCount, 1), len(entries))

	jobs := make(chan mediaInventoryEntry)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	verify := func(entry mediaInventoryEntry) error {
		data, ok := e.parts.Get(entry.partPath)
		if !ok {
			return fmt.Errorf("media inventory part missing: %s", entry.partPath)
		}
		sum := sha256.Sum256(data)
		got := hex.EncodeToString(sum[:])
		if got != entry.hash {
			return fmt.Errorf("media checksum mismatch for %s", entry.partPath)
		}
		return nil
	}

	startWorker := func() {
		defer wg.Done()
		for entry := range jobs {
			if err := verify(entry); err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
		}
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go startWorker()
	}

	for _, entry := range entries {
		select {
		case jobs <- entry:
		case err := <-errCh:
			close(jobs)
			wg.Wait()
			return err
		}
	}
	close(jobs)
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (e *PresentationEditor) snapshotMediaInventoryEntries() []mediaInventoryEntry {
	e.mediaMu.Lock()
	defer e.mediaMu.Unlock()

	entries := make([]mediaInventoryEntry, 0, len(e.mediaInventory))
	for hash, partPath := range e.mediaInventory {
		entries = append(entries, mediaInventoryEntry{hash: hash, partPath: partPath})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].hash < entries[j].hash
	})
	return entries
}
