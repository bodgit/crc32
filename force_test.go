package crc32_test

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/bodgit/crc32"
)

func ExampleForceCRC32() {
	f, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}

		if err := os.RemoveAll(f.Name()); err != nil {
			panic(err)
		}
	}()

	b := []byte("This is a test 0000")

	if _, err := f.Write(b); err != nil {
		panic(err)
	}

	if err := crc32.ForceCRC32(f, int64(len(b)-4), 0xdeadbeef); err != nil {
		panic(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	if _, err := io.ReadFull(f, b); err != nil {
		panic(err)
	}

	fmt.Print(hex.Dump(b))

	// Output: 00000000  54 68 69 73 20 69 73 20  61 20 74 65 73 74 20 99  |This is a test .|
	// 00000010  3c 5f 27                                          |<_'|
}
