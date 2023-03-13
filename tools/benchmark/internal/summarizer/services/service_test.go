package services

import (
	"benchmark/internal/encoder"
	"testing"
)

func TestParseSolutionLogName(t *testing.T) {
	encoder_, function, step, targetHash, err := parseSolutionLogName("transalg_md4_41_ffffffffffffffffffffffffffffffff_dobbertin31.cnf.cadical_c1000000.cnf.march_n2460.cubes.cube79837.kissat.log")
	if encoder_ != encoder.Transalg {
		t.Errorf("Expected transalg, got %s", encoder_)
	}
	if function != encoder.Md4 {
		t.Errorf("Expected md4, got %s", function)
	}
	if step != 41 {
		t.Errorf("Expected 41, got %d", step)
	}
	if targetHash != "ffffffffffffffffffffffffffffffff" {
		t.Errorf("Expected ffffffffffffffffffffffffffffffff, got %s", targetHash)
	}
	if err != nil {
		t.Errorf("Got error %s", err)
	}
}

func TestParseSolutionLogName2(t *testing.T) {
	encoder_, function, step, targetHash, err := parseSolutionLogName("saeed_e_md5_35_00000000000000000000000000000000_dobbertin31.cnf.cadical_c1000000.cnf.march_n2460.cubes.cube79837.kissat.log")
	if encoder_ != encoder.SaeedE {
		t.Errorf("Expected saeed_e, got %s", encoder_)
	}
	if function != encoder.Md5 {
		t.Errorf("Expected md5, got %s", function)
	}
	if step != 35 {
		t.Errorf("Expected 35, got %d", step)
	}
	if targetHash != "00000000000000000000000000000000" {
		t.Errorf("Expected 00000000000000000000000000000000, got %s", targetHash)
	}
	if err != nil {
		t.Errorf("Got error %s", err)
	}

}
func TestParseSolutionLogName3(t *testing.T) {
	encoder_, function, step, targetHash, err := parseSolutionLogName("transalg_md5_28_ffffffffffffffffffffffffffffffff_dobbertin0.cnf.cadical_c100.cnf.march_n6640.cubes.cube67217.kissat")
	if encoder_ != encoder.Transalg {
		t.Errorf("Expected transalg, got %s", encoder_)
	}
	if function != encoder.Md5 {
		t.Errorf("Expected md5, got %s", function)
	}
	if step != 28 {
		t.Errorf("Expected 28, got %d", step)
	}
	if targetHash != "ffffffffffffffffffffffffffffffff" {
		t.Errorf("Expected ffffffffffffffffffffffffffffffff, got %s", targetHash)
	}
	if err != nil {
		t.Errorf("Got error %s", err)
	}
}
