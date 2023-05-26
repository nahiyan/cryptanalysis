package services

import (
	"encoding/hex"
	"testing"
)

func TestSha256Case1(t *testing.T) {
	sha256Svc := Sha256Service{}
	bytes, _ := hex.DecodeString("e57d8668a57d8668a57d8668bc8c857ba57d8668a57d8668a57d8668cb0a1178a57d8668a57d8668a57d8668307bc4e7ad02e703e1516b23981c2a75c08ea9f7")
	digest, err := sha256Svc.Run(bytes, 64, false)
	if err != nil {
		t.Fatal("failed to compute MD5 hash: ", err)
	}

	if digest != "73703799050cdca89b07dd7605078b469b6e3608d6363ddf8d97bffe009dbc6f" {
		t.Fatalf("got hash = %s but expected 73703799050cdca89b07dd7605078b469b6e3608d6363ddf8d97bffe009dbc6f\n", digest)
	}
}

// func TestSha256Case2(t *testing.T) {
// 	sha256Svc := Sha256Service{}
// 	bytes, _ := hex.DecodeString("dfe6feebc4437a8511af5182e3b13f035103e1fcea231da2c3b513d1b95fa9d77a2a331c2ddf2607699a2daec1827561fe80aeedcf45b09a5b596c8fd0265347")
// 	digest, err := sha256Svc.Run(bytes, 28, false)
// 	if err != nil {
// 		t.Fatal("failed to compute MD5 hash: ", err)
// 	}

// 	if digest != "ffffffffffffffffffffffffffffffff" {
// 		t.Fatalf("got hash = %s but expected all-one-bit hash\n", digest)
// 	}
// }
