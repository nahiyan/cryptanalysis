package services

import (
	"encoding/hex"
	"testing"
)

func Test40Md4N1(t *testing.T) {
	md4Svc := Md4Service{}
	bytes, _ := hex.DecodeString("e57d8668a57d8668a57d8668bc8c857ba57d8668a57d8668a57d8668cb0a1178a57d8668a57d8668a57d8668307bc4e7ad02e703e1516b23981c2a75c08ea9f7")
	digest, err := md4Svc.Run(bytes, 40, false)
	if err != nil {
		t.Fatal("failed to compute MD4 hash: ", err)
	}

	if digest != "00000000000000000000000000000000" {
		t.Fatalf("got hash = %s but expected all-zero-bit hash\n", digest)
	}
}

func Test40Md4N2(t *testing.T) {
	md4Svc := Md4Service{}
	bytes, _ := hex.DecodeString("e57d8668a57d8668a57d86681d236482a57d8668a57d8668a57d866897a13204a57d8668a57d8668a57d86680991ede3301e2ac35bed2a3de167a833890d22f0")
	digest, err := md4Svc.Run(bytes, 40, false)
	if err != nil {
		t.Fatal("failed to compute MD4 hash: ", err)
	}

	if digest != "ffffffffffffffffffffffffffffffff" {
		t.Errorf("got hash = %s but expected all-one-bit hash\n", digest)
	}
}

func Test43Md4N1(t *testing.T) {
	md4Svc := Md4Service{}
	bytes, _ := hex.DecodeString("a57d8668a57d8668a57d8668f48a97a3a57d8668a57d8668a57d8668d330e8eda57d8668a57d8668a57d866837c9ca21e1df551f7f49d66a135a1c939e744bdb")
	digest, err := md4Svc.Run(bytes, 43, false)
	if err != nil {
		t.Fatal("failed to compute MD4 hash: ", err)
	}

	if digest != "00000000000000000000000000000000" {
		t.Fatalf("got hash = %s but expected all-zero-bit hash\n", digest)
	}
}

func Test43Md4N2(t *testing.T) {
	md4Svc := Md4Service{}
	bytes, _ := hex.DecodeString("a57d8668a57d8668a57d8668b289afa0a57d8668a57d8668a57d8668af2c850ea57d8668a57d8668a57d866819c5ce09cae6b29eb2595b20ab3a433df6cdee42")
	digest, err := md4Svc.Run(bytes, 43, false)
	if err != nil {
		t.Fatal("failed to compute MD4 hash: ", err)
	}

	if digest != "00000000000000000000000000000000" {
		t.Fatalf("got hash = %s but expected all-zero-bit hash\n", digest)
	}
}

func Test43Md4N3(t *testing.T) {
	md4Svc := Md4Service{}
	bytes, _ := hex.DecodeString("a57d8668a57d8668a57d866882ef987aa57d8668a57d8668a57d8668e18fbc3ba57d8668a57d8668a57d8668558f3513bf09004d8fb490dd0502eca9bd0e1a80")
	digest, err := md4Svc.Run(bytes, 43, false)
	if err != nil {
		t.Fatal("failed to compute MD4 hash: ", err)
	}

	if digest != "ffffffffffffffffffffffffffffffff" {
		t.Errorf("got hash = %s but expected all-one-bit hash\n", digest)
	}
}

func Test35MD4N1(t *testing.T) {
	md4Svc := Md4Service{}
	bytes, _ := hex.DecodeString("f082d2e8a57d8668a57d86689c1cdcd8a57d8668a57d8668a57d8668a77dc6e8a57d8668a57d8668a57d8668a3739b2fbd94a012ba83a52659d8d8d948030a00")
	digest, err := md4Svc.Run(bytes, 35, false)
	if err != nil {
		t.Fatal("failed to compute MD4 hash: ", err)
	}

	if digest != "ffffffffffffffffffffffffffffffff" {
		t.Errorf("got hash = %s but expected all-one-bit hash\n", digest)
	}
}

func Test31MD4(t *testing.T) {
	md4Svc := Md4Service{}
	bytes, _ := hex.DecodeString("a57d8668a57d8668a57d86687894e208a57d8668a57d8668a57d86685c4116c6a57d8668a57d8668a57d86688d764bfb9bcadf0a5fc20ab5930c12bd95c354d5")
	digest, err := md4Svc.Run(bytes, 31, true)
	if err != nil {
		t.Fatal("failed to compute MD4 hash: ", err)
	}

	if digest != "ffffffffffffffffffffffffffffffff" {
		t.Errorf("got hash = %s but expected all-one-bit hash\n", digest)
	}
}

func TestBytesToUint32(t *testing.T) {
	bytes, _ := hex.DecodeString("67452301efcdab89")
	result := toUint32Slice(bytes)[0]
	if result != 1732584193 {
		t.Fatalf("got 1st word = %d but expected 1732584193\n", result)
	}

	result2 := toUint32Slice(bytes)[1]
	if result2 != 4023233417 {
		t.Fatalf("got 1st word = %d but expected 4023233417\n", result2)
	}
}
