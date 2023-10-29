package services

import (
	"encoding/hex"
	"testing"
)

func Test40Md5N1(t *testing.T) {
	md5Svc := Md5Service{}
	bytes, _ := hex.DecodeString("d825e4fba73fcaa9660cd53db93085154677d4e0cadcee6240722cb3f41a4b12ac2fdec39cbcb4a3ffcca30f9a0e2026475763e530ce233bbef0cd5701a6b39d")
	digest, err := md5Svc.Run(bytes, 28, false)
	if err != nil {
		t.Fatal("failed to compute MD5 hash: ", err)
	}

	if digest != "00000000000000000000000000000000" {
		t.Fatalf("got hash = %s but expected all-zero-bit hash\n", digest)
	}
}

func Test40Md5N2(t *testing.T) {
	md5Svc := Md5Service{}
	bytes, _ := hex.DecodeString("dfe6feebc4437a8511af5182e3b13f035103e1fcea231da2c3b513d1b95fa9d77a2a331c2ddf2607699a2daec1827561fe80aeedcf45b09a5b596c8fd0265347")
	digest, err := md5Svc.Run(bytes, 28, false)
	if err != nil {
		t.Fatal("failed to compute MD5 hash: ", err)
	}

	if digest != "ffffffffffffffffffffffffffffffff" {
		t.Fatalf("got hash = %s but expected all-one-bit hash\n", digest)
	}
}
