package awsconfig

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"syscall"

	"gopkg.in/ini.v1"
)

// lockedFile holds an open, locked file handle and its parsed ini content.
type lockedFile struct {
	file *os.File
	cfg  *ini.File
}

// lockAndLoad opens the credentials file, acquires an exclusive lock, reads it into memory, and parses it.
func lockAndLoad(path string) (*lockedFile, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to lock %s: %w", path, err)
	}

	// Read content into memory so we can parse without closing the fd
	content, err := io.ReadAll(f)
	if err != nil {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var cfg *ini.File
	if len(bytes.TrimSpace(content)) == 0 {
		cfg = ini.Empty()
	} else {
		cfg, err = ini.Load(content)
		if err != nil {
			cfg = ini.Empty()
		}
	}

	return &lockedFile{file: f, cfg: cfg}, nil
}

// save writes the ini content back to the file and truncates any leftover bytes.
func (lf *lockedFile) save(path string) error {
	if _, err := lf.file.Seek(0, 0); err != nil {
		return err
	}
	if err := lf.file.Truncate(0); err != nil {
		return err
	}
	_, err := lf.cfg.WriteTo(lf.file)
	return err
}

// unlock releases the file lock and closes the file.
func (lf *lockedFile) unlock() {
	_ = syscall.Flock(int(lf.file.Fd()), syscall.LOCK_UN)
	lf.file.Close()
}
