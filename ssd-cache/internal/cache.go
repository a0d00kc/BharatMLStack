package internal

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/Meesho/BharatMLStack/ssd-cache/internal/allocator"
	"github.com/Meesho/BharatMLStack/ssd-cache/internal/index"
	"github.com/Meesho/BharatMLStack/ssd-cache/internal/memtable"
)

const (
	// File flags for testing
	O_DIRECT  = 0x4000
	O_WRONLY  = syscall.O_WRONLY
	O_APPEND  = syscall.O_APPEND
	O_CREAT   = syscall.O_CREAT
	FILE_MODE = 0644
)

type Cache struct {
	memtable *memtable.Memtable
	index    *index.Index
}

func NewCache(capacity int64) *Cache {
	allocator := allocator.NewAlignedPageAllocator(allocator.AlignedPageAllocatorConfig{
		PageSizeAlignement: 4096,
		Multiplier:         int(capacity / 4096),
		MaxPages:           1000,
	})

	// Create temporary directory
	tmpDir := os.TempDir()
	filename := filepath.Join(tmpDir, "test_memtable.dat")

	// Open file with DIRECT_IO, WRITE_ONLY, APPEND_ONLY flags
	flags := O_DIRECT | O_WRONLY | O_APPEND | O_CREAT
	fd, err := syscall.Open(filename, flags, FILE_MODE)
	if err != nil {
		// If DIRECT_IO is not supported, fall back to regular flags
		log.Println("DIRECT_IO not supported, falling back to regular flags: %v", err)
		flags = O_WRONLY | O_APPEND | O_CREAT
		fd, err = syscall.Open(filename, flags, FILE_MODE)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
	}

	file := os.NewFile(uintptr(fd), filename)
	if file == nil {
		log.Fatalf("Failed to create file from fd")
	}

	memtable := memtable.NewMemtableV2(file, 0, capacity, allocator)
	index := index.NewIndex()
	return &Cache{
		memtable: memtable,
		index:    index,
	}
}

func (c *Cache) Put(key string, value []byte) {
	offset, length := c.memtable.Put(value)
	c.index.Put(key, []byte{byte(offset), byte(length)})
}
