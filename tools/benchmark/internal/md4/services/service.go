package services

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func leftRotate(x, s uint32) uint32 {
	return (x<<s)&0xffffffff | (x >> (32 - s))
}

func sum(a, b uint32) uint32 {
	return uint32(a + b)
}

func f(x, y, z uint32) uint32 {
	return (x & y) | (^x & z)
}

func g(x, y, z uint32) uint32 {
	return (x & y) | (x & z) | (y & z)
}

func h(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func ff(a, b, c, d, m, s uint32) uint32 {
	a = sum(sum(a, f(b, c, d)), m)
	return leftRotate(a, s)
}

func gg(a, b, c, d, m, s uint32) uint32 {
	a = sum(sum(sum(a, g(b, c, d)), m), 0x5a827999)
	return leftRotate(a, s)
}

func hh(a, b, c, d, m, s uint32) uint32 {
	a = sum(sum(sum(a, h(b, c, d)), m), 0x6ed9eba1)
	return leftRotate(a, s)
}

func toUint32Slice(bytes []byte) []uint32 {
	words := []uint32{}
	for i := 0; i < len(bytes); i += 4 {
		word := binary.BigEndian.Uint32(bytes[i:])
		words = append(words, word)
	}

	return words
}

func (md4Svc *Md4Service) Run(message_ []byte, steps int) (string, error) {
	digest := ""

	if len(message_) != 64 {
		return digest, errors.New("message must be exactly 512 bits long")
	}
	message := toUint32Slice(message_)

	// Initial
	var a_ uint32 = 0x67452301
	var b_ uint32 = 0xefcdab89
	var c_ uint32 = 0x98badcfe
	var d_ uint32 = 0x10325476

	var a uint32 = a_
	var b uint32 = b_
	var c uint32 = c_
	var d uint32 = d_
	h := []uint32{a, b, c, d}

	s_ := []uint32{3, 7, 11, 19, 3, 5, 9, 13, 3, 9, 11, 15}

	for step := 1; step <= steps; step++ {
		if step <= 16 {
			// Round 1
			i := (16 - (step - 1)) % 4
			s := s_[0:]
			k := step - 1
			h[i] = ff(h[i], h[(i+1)%4], h[(i+2)%4], h[(i+3)%4], message[k], s[(step-1)%4])
		} else if step <= 32 {
			// Round 2
			i := (32 - (step - 1)) % 4
			s := s_[4:]
			k := []uint32{0, 4, 8, 12, 1, 5, 9, 13, 2, 6, 10, 14, 3, 7, 11, 15}
			h[i] = gg(h[i], h[(i+1)%4], h[(i+2)%4], h[(i+3)%4], message[k[step-1-16]], s[(step-1)%4])
		} else {
			// Round 3
			i := (48 - (step - 1)) % 4
			s := s_[8:]
			k := []uint32{0, 8, 4, 12, 2, 10, 6, 14, 1, 9, 5, 13,
				3, 11, 7, 15}
			h[i] = hh(h[i], h[(i+1)%4], h[(i+2)%4], h[(i+3)%4], message[k[step-1-32]], s[(step-1-32)%4])
		}
	}

	// Final
	a = h[0]
	b = h[1]
	c = h[2]
	d = h[3]
	// a += a_
	// b += b_
	// c += c_
	// d += d_
	// fmt.Println("Final", a, b, c, d)
	digest_ := make([]byte, 16)
	binary.LittleEndian.PutUint32(digest_, a)
	binary.LittleEndian.PutUint32(digest_[4:], b)
	binary.LittleEndian.PutUint32(digest_[8:], c)
	binary.LittleEndian.PutUint32(digest_[12:], d)
	digest = fmt.Sprintf("%x", digest_)

	return digest, nil
}
