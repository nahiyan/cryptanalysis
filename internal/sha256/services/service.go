package services

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func rightRotate(x, s uint32) uint32 {
	return (x >> s) | (x << (32 - s))
}

func ch(x, y, z uint32) uint32 {
	return (x & y) ^ (^x & z)
}

func maj(x, y, z uint32) uint32 {
	return (x & y) ^ (x & z) ^ (y & z)
}

func ep0(x uint32) uint32 {
	return rightRotate(x, 2) ^ rightRotate(x, 13) ^ rightRotate(x, 22)
}

func ep1(x uint32) uint32 {
	return rightRotate(x, 6) ^ rightRotate(x, 11) ^ rightRotate(x, 25)
}

func sig0(x uint32) uint32 {
	return rightRotate(x, 7) ^ rightRotate(x, 18) ^ (x >> 3)
}

func sig1(x uint32) uint32 {
	return rightRotate(x, 17) ^ rightRotate(x, 19) ^ (x >> 10)
}

func (sha256Svc *Sha256Service) Run(message []byte, steps int, addChainingVars bool) (string, error) {
	k := []uint32{
		0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5, 0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174, 0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da, 0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967, 0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85, 0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070, 0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3, 0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
	}

	if len(message) != 64 {
		return "", errors.New("message must be exactly 512 bits long")
	}

	// Initial
	var h0 uint32 = 0x6a09e667
	var h1 uint32 = 0xbb67ae85
	var h2 uint32 = 0x3c6ef372
	var h3 uint32 = 0xa54ff53a
	var h5 uint32 = 0x9b05688c
	var h4 uint32 = 0x510e527f
	var h6 uint32 = 0x1f83d9ab
	var h7 uint32 = 0x5be0cd19

	// Message expansion
	expandedMessage := make([]uint32, 64)
	for i := 0; i < 16; i++ {
		expandedMessage[i] = binary.BigEndian.Uint32(message[i*4:])
		// fmt.Printf("W[%d] = %08x\n", i, expandedMessage[i])
	}
	for i := 16; i < steps; i++ {
		expandedMessage[i] = sig1(expandedMessage[i-2]) + expandedMessage[i-7] + sig0(expandedMessage[i-15]) + expandedMessage[i-16]
		// fmt.Printf("W[%d] = %08x\n", i, expandedMessage[i])
	}

	var a uint32 = h0
	var b uint32 = h1
	var c uint32 = h2
	var d uint32 = h3
	var e uint32 = h4
	var f uint32 = h5
	var g uint32 = h6
	var h uint32 = h7
	var t1, t2 uint32

	for i := 0; i < steps; i++ {
		t1 = h + ep1(e) + ch(e, f, g) + k[i] + expandedMessage[i]
		t2 = ep0(a) + maj(a, b, c)
		h = g
		g = f
		f = e
		e = d + t1
		d = c
		c = b
		b = a
		a = t1 + t2
		// fmt.Printf("Step %d %08x %08x\n", i, a, e)
	}

	h0 = a
	h1 = b
	h2 = c
	h3 = d
	h4 = e
	h5 = f
	h6 = g
	h7 = h

	if addChainingVars {
		h0 += a
		h1 += b
		h2 += c
		h3 += d
		h4 += e
		h5 += f
		h6 += g
		h7 += h
	}

	digest := fmt.Sprintf("%08x%08x%08x%08x%08x%08x%08x%08x", h0, h1, h2, h3, h4, h5, h6, h7)

	return digest, nil
}
