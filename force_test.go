package crc32

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func ExampleForceCRC32() {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Close()
		os.RemoveAll(f.Name())
	}()

	b := []byte("This is a test 0000")

	if n, err := f.Write(b); err != nil || n != len(b) {
		if err != nil {
			panic(err)
		}
		panic(errors.New("short write"))
	}

	if err := ForceCRC32(f, int64(len(b)-4), 0xdeadbeef); err != nil {
		panic(err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	if n, err := f.Read(b); err != nil || n != len(b) {
		if err != nil {
			panic(err)
		}
		panic(errors.New("short read"))
	}

	fmt.Print(hex.Dump(b))

	// Output: 00000000  54 68 69 73 20 69 73 20  61 20 74 65 73 74 20 99  |This is a test .|
	// 00000010  3c 5f 27                                          |<_'|

}
