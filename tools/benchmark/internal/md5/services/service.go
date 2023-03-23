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
	return (x & z) | (y & ^z)
}

func h(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func i(x, y, z uint32) uint32 {
	return y ^ (x | ^z)
}

func ff(a, b, c, d, m, s, t uint32) uint32 {
	a = sum(sum(sum(a, f(b, c, d)), m), t)
	return b + leftRotate(a, s)
}

func gg(a, b, c, d, m, s, t uint32) uint32 {
	a = sum(sum(sum(a, g(b, c, d)), m), t)
	return b + leftRotate(a, s)
}

func hh(a, b, c, d, m, s, t uint32) uint32 {
	a = sum(sum(sum(a, h(b, c, d)), m), t)
	return b + leftRotate(a, s)
}

func ii(a, b, c, d, m, s, t uint32) uint32 {
	a = sum(sum(sum(a, i(b, c, d)), m), t)
	return b + leftRotate(a, s)
}

func toUint32Slice(bytes []byte) []uint32 {
	words := []uint32{}
	for i := 0; i < len(bytes); i += 4 {
		word := binary.BigEndian.Uint32(bytes[i:])
		words = append(words, word)
	}

	return words
}

func (md4Svc *Md5Service) Run(message_ []byte, steps int, addChainingVars bool) (string, error) {
	if len(message_) != 64 {
		return "", errors.New("message must be exactly 512 bits long")
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

	// Round 1: ff(a,b,c,d,M_j,s,t_j) 1-16

	if steps >= 1 {
		a = ff(a, b, c, d, message[0], 7, 0xd76aa478)
	}
	if steps >= 2 {
		d = ff(d, a, b, c, message[1], 12, 0xe8c7b756)
	}
	if steps >= 3 {
		c = ff(c, d, a, b, message[2], 17, 0x242070db)
	}
	if steps >= 4 {
		b = ff(b, c, d, a, message[3], 22, 0xc1bdceee)
	}

	if steps >= 5 {
		a = ff(a, b, c, d, message[4], 7, 0xf57c0faf)
	}
	if steps >= 6 {
		d = ff(d, a, b, c, message[5], 12, 0x4787c62a)
	}
	if steps >= 7 {
		c = ff(c, d, a, b, message[6], 17, 0xa8304613)
	}
	if steps >= 8 {
		b = ff(b, c, d, a, message[7], 22, 0xfd469501)
	}

	if steps >= 9 {
		a = ff(a, b, c, d, message[8], 7, 0x698098d8)
	}
	if steps >= 10 {
		d = ff(d, a, b, c, message[9], 12, 0x8b44f7af)
	}
	if steps >= 11 {
		c = ff(c, d, a, b, message[10], 17, 0xffff5bb1)
	}
	if steps >= 12 {
		b = ff(b, c, d, a, message[11], 22, 0x895cd7be)
	}

	if steps >= 13 {
		a = ff(a, b, c, d, message[12], 7, 0x6b901122)
	}
	if steps >= 14 {
		d = ff(d, a, b, c, message[13], 12, 0xfd987193)
	}
	if steps >= 15 {
		c = ff(c, d, a, b, message[14], 17, 0xa679438e)
	}
	if steps >= 16 {
		b = ff(b, c, d, a, message[15], 22, 0x49b40821)
	}

	// Round 2: gg(a,b,c,d,M_j,s,t_j) 17-32

	if steps >= 17 {
		a = gg(a, b, c, d, message[1], 5, 0xf61e2562)
	}
	if steps >= 18 {
		d = gg(d, a, b, c, message[6], 9, 0xc040b340)
	}
	if steps >= 19 {
		c = gg(c, d, a, b, message[11], 14, 0x265e5a51)
	}
	if steps >= 20 {
		b = gg(b, c, d, a, message[0], 20, 0xe9b6c7aa)
	}

	if steps >= 21 {
		a = gg(a, b, c, d, message[5], 5, 0xd62f105d)
	}
	if steps >= 22 {
		d = gg(d, a, b, c, message[10], 9, 0x02441453)
	}
	if steps >= 23 {
		c = gg(c, d, a, b, message[15], 14, 0xd8a1e681)
	}
	if steps >= 24 {
		b = gg(b, c, d, a, message[4], 20, 0xe7d3fbc8)
	}

	if steps >= 25 {
		a = gg(a, b, c, d, message[9], 5, 0x21e1cde6)
	}
	if steps >= 26 {
		d = gg(d, a, b, c, message[14], 9, 0xc33707d6)
	}
	if steps >= 27 {
		c = gg(c, d, a, b, message[3], 14, 0xf4d50d87)
	}
	if steps >= 28 {
		b = gg(b, c, d, a, message[8], 20, 0x455a14ed)
	}

	if steps >= 29 {
		a = gg(a, b, c, d, message[13], 5, 0xa9e3e905)
	}
	if steps >= 30 {
		d = gg(d, a, b, c, message[2], 9, 0xfcefa3f8)
	}
	if steps >= 31 {
		c = gg(c, d, a, b, message[7], 14, 0x676f02d9)
	}
	if steps >= 32 {
		b = gg(b, c, d, a, message[12], 20, 0x8d2a4c8a)
	}

	// Round 3: hh(a,b,c,d,M_j,s,t_j) 33-48

	if steps >= 33 {
		a = hh(a, b, c, d, message[5], 4, 0xfffa3942)
	}
	if steps >= 34 {
		d = hh(d, a, b, c, message[8], 11, 0x8771f681)
	}
	if steps >= 35 {
		c = hh(c, d, a, b, message[11], 16, 0x6d9d6122)
	}
	if steps >= 36 {
		b = hh(b, c, d, a, message[14], 23, 0xfde5380c)
	}

	if steps >= 37 {
		a = hh(a, b, c, d, message[1], 4, 0xa4beea44)
	}
	if steps >= 38 {
		d = hh(d, a, b, c, message[4], 11, 0x4bdecfa9)
	}
	if steps >= 39 {
		c = hh(c, d, a, b, message[7], 16, 0xf6bb4b60)
	}
	if steps >= 40 {
		b = hh(b, c, d, a, message[10], 23, 0xbebfbc70)
	}

	if steps >= 41 {
		a = hh(a, b, c, d, message[13], 4, 0x289b7ec6)
	}
	if steps >= 42 {
		d = hh(d, a, b, c, message[0], 11, 0xeaa127fa)
	}
	if steps >= 43 {
		c = hh(c, d, a, b, message[3], 16, 0xd4ef3085)
	}
	if steps >= 44 {
		b = hh(b, c, d, a, message[6], 23, 0x04881d05)
	}

	if steps >= 45 {
		a = hh(a, b, c, d, message[9], 4, 0xd9d4d039)
	}
	if steps >= 46 {
		d = hh(d, a, b, c, message[12], 11, 0xe6db99e5)
	}
	if steps >= 47 {
		c = hh(c, d, a, b, message[15], 16, 0x1fa27cf8)
	}
	if steps >= 48 {
		b = hh(b, c, d, a, message[2], 23, 0xc4ac5665)
	}

	// Round 4: ii(a,b,c,d,M_j,s,t_j) 49-64

	if steps >= 49 {
		a = ii(a, b, c, d, message[0], 6, 0xf4292244)
	}
	if steps >= 50 {
		d = ii(d, a, b, c, message[7], 10, 0x432aff97)
	}
	if steps >= 51 {
		c = ii(c, d, a, b, message[14], 15, 0xab9423a7)
	}
	if steps >= 52 {
		b = ii(b, c, d, a, message[5], 21, 0xfc93a039)
	}

	if steps >= 53 {
		a = ii(a, b, c, d, message[12], 6, 0x655b59c3)
	}
	if steps >= 54 {
		d = ii(d, a, b, c, message[3], 10, 0x8f0ccc92)
	}
	if steps >= 55 {
		c = ii(c, d, a, b, message[10], 15, 0xffeff47d)
	}
	if steps >= 56 {
		b = ii(b, c, d, a, message[1], 21, 0x85845dd1)
	}

	if steps >= 57 {
		a = ii(a, b, c, d, message[8], 6, 0x6fa87e4f)
	}
	if steps >= 58 {
		d = ii(d, a, b, c, message[15], 10, 0xfe2ce6e0)
	}
	if steps >= 59 {
		c = ii(c, d, a, b, message[6], 15, 0xa3014314)
	}
	if steps >= 60 {
		b = ii(b, c, d, a, message[13], 21, 0x4e0811a1)
	}

	if steps >= 61 {
		a = ii(a, b, c, d, message[4], 6, 0xf7537e82)
	}
	if steps >= 62 {
		d = ii(d, a, b, c, message[11], 10, 0xbd3af235)
	}
	if steps >= 63 {
		c = ii(c, d, a, b, message[2], 15, 0x2ad7d2bb)
	}
	if steps >= 64 {
		b = ii(b, c, d, a, message[9], 21, 0xeb86d391)
	}

	// Add chaining variables
	if addChainingVars {
		a += a_
		b += b_
		c += c_
		d += d_
	}

	digest_ := make([]byte, 16)
	binary.BigEndian.PutUint32(digest_, a)
	binary.BigEndian.PutUint32(digest_[4:], b)
	binary.BigEndian.PutUint32(digest_[8:], c)
	binary.BigEndian.PutUint32(digest_[12:], d)
	digest := fmt.Sprintf("%x", digest_)

	return digest, nil
}
