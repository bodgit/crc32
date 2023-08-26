/*
Package crc32 is an implementation of an algorithm to modify a file so that its
CRC-32 checksum matches a given value. This requires four sacrificial bytes in
the file that will be modified to generate the desired value.
*/
package crc32

import (
	"errors"
	"hash/crc32"
	"io"
	"math/bits"
)

// Based on the code at https://www.nayuki.io/page/forcing-a-files-crc-to-any-value

const polynomial = 0x104c11db7

func multiplyMod(x, y uint) (z uint) {
	for z = 0; y != 0; {
		z ^= x * (y & 1)
		y >>= 1
		x <<= 1

		if (x>>32)&1 != 0 {
			x ^= polynomial
		}
	}

	return
}

func powMod(x, y uint) (z uint) {
	for z = 1; y != 0; {
		if y&1 != 0 {
			z = multiplyMod(z, x)
		}

		x = multiplyMod(x, x)
		y >>= 1
	}

	return
}

func divideAndRemainder(x, y uint) (uint, uint, error) {
	if y == 0 {
		return 0, 0, errors.New("divide by zero")
	}

	if x == 0 {
		return 0, 0, nil
	}

	ydeg, z := getDegree(y), uint(0)
	for i := getDegree(x) - ydeg; i >= 0; i-- {
		if (x>>(i+ydeg))&1 != 0 {
			x ^= y << i
			z |= 1 << i
		}
	}

	return z, x, nil
}

func reciprocalMod(x uint) (uint, error) {
	y := x
	x = polynomial
	a, b := uint(0), uint(1)

	for y != 0 {
		q, r, err := divideAndRemainder(x, y)
		if err != nil {
			return 0, err
		}

		c := a ^ multiplyMod(q, b)
		x, y, a, b = y, r, b, c
	}

	if x == 1 {
		return a, nil
	}

	return 0, errors.New("reciprocal does not exist")
}

func getDegree(x uint) int {
	return bits.Len(x) - 1
}

func generateCRC32(rs io.ReadSeeker) (uint32, error) {
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	h := crc32.NewIEEE()

	if _, err := io.Copy(h, rs); err != nil {
		return 0, err
	}

	crc := h.Sum(nil)

	return uint32(crc[0])<<24 | uint32(crc[1])<<16 | uint32(crc[2])<<8 | uint32(crc[3]), nil
}

//nolint:cyclop
func overwriteCRC32(rws io.ReadWriteSeeker, length, offset int64, currentCRC, desiredCRC uint32) error {
	a, err := reciprocalMod(powMod(2, uint((length-offset)*8)))
	if err != nil {
		return err
	}

	delta := uint32(multiplyMod(a, uint(bits.Reverse32(currentCRC^desiredCRC))))

	if _, err = rws.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	b := make([]byte, 4)
	if n, err := rws.Read(b); err != nil || n != len(b) {
		if err != nil {
			return err
		}

		return errors.New("short read")
	}

	for i := 0; i < len(b); i++ {
		b[i] ^= byte((bits.Reverse32(delta) >> (i * 8)) & 0xff)
	}

	if _, err = rws.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if n, err := rws.Write(b); err != nil || n != len(b) {
		if err != nil {
			return err
		}

		return errors.New("short write")
	}

	return nil
}

// ForceCRC32 takes an io.ReadWriteSeeker and updates the four bytes starting
// at offset such that the IEEE CRC-32 value of the whole stream matches
// desiredCRC. If the CRC-32 value already matches, no changes are made.
func ForceCRC32(rws io.ReadWriteSeeker, offset int64, desiredCRC uint32) error {
	length, err := rws.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	if offset < 0 || offset+4 > length {
		return errors.New("invalid byte offset")
	}

	currentCRC, err := generateCRC32(rws)
	if err != nil {
		return err
	}

	if currentCRC == desiredCRC {
		return nil
	}

	if err := overwriteCRC32(rws, length, offset, currentCRC, desiredCRC); err != nil {
		return err
	}

	newCRC, err := generateCRC32(rws)
	if err != nil {
		return err
	}

	if newCRC != desiredCRC {
		return errors.New("new CRC does not match desired CRC")
	}

	return nil
}
