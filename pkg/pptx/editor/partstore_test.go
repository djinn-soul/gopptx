package editor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestPartStoreGetSetDelete(t *testing.T) {
	ps := newPartStoreFromMap(map[string][]byte{
		"a.xml": []byte("alpha"),
		"b.xml": []byte("bravo"),
	})

	data, ok := ps.Get("a.xml")
	if !ok || string(data) != "alpha" {
		t.Fatalf("expected alpha, got %q ok=%v", data, ok)
	}

	ps.Set("c.xml", []byte("charlie"))
	data, ok = ps.Get("c.xml")
	if !ok || string(data) != "charlie" {
		t.Fatalf("expected charlie after Set")
	}

	ps.Delete("a.xml")
	if ps.Has("a.xml") {
		t.Fatalf("expected a.xml to be deleted")
	}
	_, ok = ps.Get("a.xml")
	if ok {
		t.Fatalf("expected Get to return false for deleted key")
	}
}

func TestPartStoreModifiedOverridesOriginal(t *testing.T) {
	ps := newPartStoreFromMap(map[string][]byte{
		"file.xml": []byte("original"),
	})

	ps.Set("file.xml", []byte("modified"))
	data, ok := ps.Get("file.xml")
	if !ok || string(data) != "modified" {
		t.Fatalf("expected modified to override original, got %q", data)
	}
}

func TestPartStoreKeys(t *testing.T) {
	ps := newPartStoreFromMap(map[string][]byte{
		"b.xml": []byte("b"),
		"a.xml": []byte("a"),
	})
	ps.Set("c.xml", []byte("c"))
	ps.Delete("a.xml")

	keys := ps.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}
	if keys[0] != "b.xml" || keys[1] != "c.xml" {
		t.Fatalf("expected [b.xml, c.xml], got %v", keys)
	}
}

func TestPartStoreKeysWithPrefix(t *testing.T) {
	ps := newPartStoreFromMap(map[string][]byte{
		"ppt/slides/slide1.xml": []byte("s1"),
		"ppt/slides/slide2.xml": []byte("s2"),
		"ppt/theme/theme1.xml":  []byte("t1"),
	})

	slides := ps.KeysWithPrefix("ppt/slides/")
	if len(slides) != 2 {
		t.Fatalf("expected 2 slide keys, got %d", len(slides))
	}
}

func TestPartStoreSnapshot(t *testing.T) {
	ps := newPartStoreFromMap(map[string][]byte{
		"a.xml": []byte("alpha"),
	})
	ps.Set("b.xml", []byte("bravo"))

	snap := ps.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries in snapshot, got %d", len(snap))
	}
	if string(snap["a.xml"]) != "alpha" || string(snap["b.xml"]) != "bravo" {
		t.Fatalf("snapshot data mismatch")
	}
}

func TestPartStoreLazyLoadFromZip(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "test.zip")

	// Create a real zip file
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	w, _ := zw.Create("part1.xml")
	_, _ = w.Write([]byte("lazy-data"))
	w2, _ := zw.Create("part2.xml")
	_, _ = w2.Write([]byte("second"))
	_ = zw.Close()
	_ = f.Close()

	// Open for lazy reading
	file, err := os.Open(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	stat, _ := file.Stat()
	zr, err := zip.NewReader(file, stat.Size())
	if err != nil {
		t.Fatal(err)
	}

	ps := newPartStoreFromZip(file, zr)
	defer func() { _ = ps.Close() }()

	// Data should not be in cache yet
	if len(ps.cache) != 0 {
		t.Fatalf("expected empty cache before Get, got %d entries", len(ps.cache))
	}

	// Has should work without loading data
	if !ps.Has("part1.xml") {
		t.Fatalf("expected Has(part1.xml) = true")
	}
	if len(ps.cache) != 0 {
		t.Fatalf("Has should not populate cache")
	}

	// Get should lazy-load
	data, ok := ps.Get("part1.xml")
	if !ok || string(data) != "lazy-data" {
		t.Fatalf("expected lazy-data, got %q", data)
	}
	if len(ps.cache) != 1 {
		t.Fatalf("expected 1 cached entry after Get, got %d", len(ps.cache))
	}

	// part2 should still not be cached
	if _, cached := ps.cache["part2.xml"]; cached {
		t.Fatalf("part2.xml should not be cached yet")
	}
}

func TestPartStoreHas(t *testing.T) {
	ps := newPartStoreFromMap(map[string][]byte{
		"exists.xml": []byte("data"),
	})

	if !ps.Has("exists.xml") {
		t.Fatalf("expected Has to return true")
	}
	if ps.Has("missing.xml") {
		t.Fatalf("expected Has to return false for missing key")
	}
}

func TestPartStoreMaterializePreservesKeysAndClosesFile(t *testing.T) {
	zipPath := buildConcurrentPartStoreZip(t, 3)
	ps := openConcurrentPartStore(t, zipPath)

	if err := ps.Materialize(); err != nil {
		t.Fatalf("materialize failed: %v", err)
	}
	if ps.file != nil {
		t.Fatalf("expected backing file to be closed after materialize")
	}
	keys := ps.KeysWithPrefix("ppt/slides/")
	if len(keys) == 0 {
		t.Fatalf("expected materialized keys to remain visible")
	}
}

func TestPartStoreCloseNilFile(t *testing.T) {
	ps := newPartStoreFromMap(map[string][]byte{})
	if err := ps.Close(); err != nil {
		t.Fatalf("Close on nil file should not error: %v", err)
	}
}

func TestPartStoreDeleteAndReAdd(t *testing.T) {
	ps := newPartStoreFromMap(map[string][]byte{
		"file.xml": []byte("original"),
	})

	ps.Delete("file.xml")
	if ps.Has("file.xml") {
		t.Fatalf("expected deleted")
	}

	ps.Set("file.xml", []byte("resurrected"))
	data, ok := ps.Get("file.xml")
	if !ok || string(data) != "resurrected" {
		t.Fatalf("expected resurrected after re-add, got %q", data)
	}
}

// Ensure newPartStoreFromMap is a true deep copy.
func TestPartStoreFromMapIsolation(t *testing.T) {
	original := map[string][]byte{
		"a.xml": []byte("alpha"),
	}
	ps := newPartStoreFromMap(original)

	// Mutating the original map should not affect the store.
	original["a.xml"] = []byte("MUTATED")

	data, _ := ps.Get("a.xml")
	if !bytes.Equal(data, []byte("alpha")) {
		t.Fatalf("PartStore should be isolated from original map mutation")
	}
}

func TestPartStoreConcurrentGet(t *testing.T) {
	zipPath := buildConcurrentPartStoreZip(t, 10)
	ps := openConcurrentPartStore(t, zipPath)
	defer func() { _ = ps.Close() }()

	names := ps.KeysWithPrefix("ppt/slides/")
	runConcurrentReadWorkers(t, ps, names, 20, 200)
}

func buildConcurrentPartStoreZip(t *testing.T, slideCount int) string {
	t.Helper()
	zipPath := filepath.Join(t.TempDir(), "concurrent.zip")
	file, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	for i := range slideCount {
		partName := filepath.ToSlash(filepath.Join("ppt", "slides", fmt.Sprintf("slide%d.xml", i)))
		entry, createErr := writer.Create(partName)
		if createErr != nil {
			t.Fatal(createErr)
		}
		if _, writeErr := entry.Write([]byte("data")); writeErr != nil {
			t.Fatal(writeErr)
		}
	}
	closeWriterErr := writer.Close()
	if closeWriterErr != nil {
		t.Fatal(closeWriterErr)
	}
	closeFileErr := file.Close()
	if closeFileErr != nil {
		t.Fatal(closeFileErr)
	}
	return zipPath
}

func openConcurrentPartStore(t *testing.T, zipPath string) *PartStore {
	t.Helper()
	file, err := os.Open(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	stat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		t.Fatal(err)
	}
	return newPartStoreFromZip(file, reader)
}

func runConcurrentReadWorkers(
	t *testing.T,
	ps *PartStore,
	names []string,
	workerCount int,
	iterations int,
) {
	t.Helper()
	var wg sync.WaitGroup
	for range workerCount {
		wg.Go(func() {
			for range iterations {
				if !assertAllPartsReadable(t, ps, names) {
					return
				}
			}
		})
	}
	wg.Wait()
}

func assertAllPartsReadable(t *testing.T, ps *PartStore, names []string) bool {
	t.Helper()
	for _, name := range names {
		if _, ok := ps.Get(name); !ok {
			t.Errorf("missing part %q during concurrent Get", name)
			return false
		}
	}
	return true
}
