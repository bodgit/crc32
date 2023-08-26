[![GitHub release](https://img.shields.io/github/v/release/bodgit/crc32)](https://github.com/bodgit/crc32/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/bodgit/crc32/build.yml?branch=main)](https://github.com/bodgit/crc32/actions?query=workflow%3ABuild)
[![Coverage Status](https://coveralls.io/repos/github/bodgit/crc32/badge.svg?branch=main)](https://coveralls.io/github/bodgit/crc32?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/bodgit/crc32)](https://goreportcard.com/report/github.com/bodgit/crc32)
[![GoDoc](https://godoc.org/github.com/bodgit/crc32?status.svg)](https://godoc.org/github.com/bodgit/crc32)
![Go version](https://img.shields.io/badge/Go-1.20-brightgreen.svg)
![Go version](https://img.shields.io/badge/Go-1.19-brightgreen.svg)

crc32
=====

An implementation of an algorithm to modify a file so that its CRC-32 checksum
matches a given value. This requires four sacrificial bytes in the file that
will be modified to generate the desired value. A small example:
```golang
f, err := os.OpenFile("somefile", O_RDWR, 0) // Remember to open read/write!
if err != nil {
        log.Fatal(err)
}
defer f.Close()

if err := crc32.ForceCRC32(f, 0, 0xdeadbeef); err != nil {
        log.Fatal(err)
}
```
In this example the first four bytes in `somefile` will be modified so that
the CRC-32 checksum of `somefile` will be `0xdeadbeef` however the bytes can
be anywhere in the file, but must be contiguous.
