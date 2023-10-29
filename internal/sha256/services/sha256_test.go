package services

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestSha256Case1(t *testing.T) {
	sha256Svc := Sha256Service{}
	message, _ := hex.DecodeString("e57d8668a57d8668a57d8668bc8c857ba57d8668a57d8668a57d8668cb0a1178a57d8668a57d8668a57d8668307bc4e7ad02e703e1516b23981c2a75c08ea9f7")
	digest, err := sha256Svc.Run(message, 64, false)
	if err != nil {
		t.Fatal("failed to compute SHA256 hash: ", err)
	}

	if digest != "73703799050cdca89b07dd7605078b469b6e3608d6363ddf8d97bffe009dbc6f" {
		t.Fatalf("got hash = %s but expected 73703799050cdca89b07dd7605078b469b6e3608d6363ddf8d97bffe009dbc6f\n", digest)
	}
}

func TestSha256Case2(t *testing.T) {
	sha256Svc := Sha256Service{}
	message, _ := hex.DecodeString("633965647a61656b66776a6b73627a306c65776b7567636a786e786c6d746879706561786e6f6b336e79386565786d7a736e6473707a7730767432656970727a")
	log.Println(message)
	digest, err := sha256Svc.Run(message, 64, true)
	if err != nil {
		t.Fatal("failed to compute SHA256 hash: ", err)
	}

	if digest != "7e24d5acc17297ecf1978c9642056f7c5bfc43114e74d426ce978a4e25973944" {
		t.Fatalf("got hash = %s but expected 7e24d5acc17297ecf1978c9642056f7c5bfc43114e74d426ce978a4e25973944\n", digest)
	}
}
func TestSha256Case3(t *testing.T) {
	sha256Svc := Sha256Service{}
	message, _ := hex.DecodeString("725a03700daa9f1b071d92dfec8282c17913134abc2eb29102d33a84278dfd290c40f8ead8bd68a00ce670c55ec7155d9f6407a8729fbfe8aa7c7c08607ae76d")
	digest, err := sha256Svc.Run(message, 27, true)
	if err != nil {
		t.Fatal("failed to compute SHA256 hash: ", err)
	}

	expectedDigest := "5864015f133494fafa42bb3594bc44f929eabb369e461e332eab27f8106467c9"
	if digest != expectedDigest {
		t.Fatalf("got hash = %s but expected %s\n", digest, expectedDigest)
	}
}
