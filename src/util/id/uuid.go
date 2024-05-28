package id

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

type uuid [16]byte

var (
	Nil uuid // empty uuid, all zeros
)

// UUID 生产uuid
func UUID() string {
	return newString()
}

func newString() string {
	return newRandom().string()
}

func (uuid uuid) string() string {
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}

func encodeHex(dst []byte, uuid uuid) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

func newRandom() uuid {
	return newRandomFromReader(rand.Reader)
}

func newRandomFromReader(r io.Reader) uuid {
	var id uuid
	_, err := io.ReadFull(r, id[:])
	if err != nil {
		return Nil
	}
	id[6] = (id[6] & 0x0f) | 0x40 // Version 4
	id[8] = (id[8] & 0x3f) | 0x80 // Variant is 10
	return id
}
