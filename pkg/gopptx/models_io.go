package gopptx

import (
	"errors"
	"os"
	"path/filepath"
)

func writeFileAtomically(path string, content []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".gopptx-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	committed := false
	defer func() {
		_ = tmpFile.Close()
		if !committed {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err = tmpFile.Write(content); err != nil {
		return err
	}
	if err = tmpFile.Chmod(perm); err != nil {
		return err
	}
	if err = tmpFile.Sync(); err != nil {
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}

	if _, err = os.Stat(path); err == nil {
		if err = os.Remove(path); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err = os.Rename(tmpPath, path); err != nil {
		return err
	}
	committed = true
	return nil
}
