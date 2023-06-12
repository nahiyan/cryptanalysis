package services

import (
	"cryptanalysis/internal/encoder"
	"testing"
)

func TestGetInstanceName(t *testing.T) {
	{
		info := encoder.InstanceInfo{
			Encoder:      encoder.Transalg,
			Steps:        35,
			Function:     "md4",
			TargetHash:   "ffffffffffffffffffffffffffffffff",
			IsXorEnabled: true,
		}
		instance := getInstanceName(info)
		expectation := "transalg_md4_35_ffffffffffffffffffffffffffffffff.cnf"
		if instance != expectation {
			t.Errorf("Expected %s but got %s instead", expectation, instance)
		}
	}

	{
		info := encoder.InstanceInfo{
			Encoder:      encoder.NejatiEncoder,
			Steps:        35,
			Function:     "md4",
			TargetHash:   "ffffffffffffffffffffffffffffffff",
			IsXorEnabled: true,
			AdderType:    encoder.Espresso,
		}
		instance := getInstanceName(info)
		expectation := "nejati_encoder_md4_35_ffffffffffffffffffffffffffffffff_espresso_xor.cnf"
		if instance != expectation {
			t.Errorf("Expected %s but got %s instead", expectation, instance)
		}
	}
}

func TestProcessInstanceName(t *testing.T) {
	{
		info, _ := processInstanceName("transalg_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.Encoder != "transalg" {
			t.Errorf("got '%s', expected 'transalg'", info.Encoder)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.Encoder != "nejati_encoder" {
			t.Errorf("got '%s', expected 'nejati_encoder'", info.Encoder)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.Function != "md4" {
			t.Errorf("got '%s', expected 'md4'", info.Function)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.Steps != 35 {
			t.Errorf("got %d, expected 35", info.Steps)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.AdderType != "espresso" {
			t.Errorf("got '%s', expected 'espresso'", info.AdderType)
		}

		info, _ = processInstanceName("nejati_encoder_md4_35_dot_matrix_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.AdderType != "dot_matrix" {
			t.Errorf("got '%s', expected 'dot_matrix'", info.AdderType)
		}

		info, _ = processInstanceName("nejati_encoder_md4_35_counter_chain_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.AdderType != "counter_chain" {
			t.Errorf("got '%s', expected 'counter_chain'", info.AdderType)
		}

		info, _ = processInstanceName("nejati_encoder_md4_35_two_operand_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.AdderType != "two_operand" {
			t.Errorf("got '%s', expected 'two_operand'", info.AdderType)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.TargetHash != "ffffffffffffffffffffffffffffffff" {
			t.Errorf("got %s, expected 'ffffffffffffffffffffffffffffffff'", info.TargetHash)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32_xor.cnf")
		dobbertin, exists := info.Dobbertin.Get()
		if !exists {
			t.Error("expected dobbertin, got none")
		}

		if dobbertin.Bits != 32 {
			t.Errorf("got %d, expected 32 dobbertin bits", dobbertin.Bits)
		}
	}

	{
		info, _ := processInstanceName("transalg_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin1_xor.cnf")
		dobbertin, exists := info.Dobbertin.Get()
		if !exists {
			t.Error("expected dobbertin, got none")
		}

		if dobbertin.Bits != 1 {
			t.Errorf("got %d, expected 1 dobbertin bits", dobbertin.Bits)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff.cnf")
		_, exists := info.Dobbertin.Get()
		if exists {
			t.Error("expected no dobbertin, got dobbertin")
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32_xor.cnf")
		if !info.IsXorEnabled {
			t.Error("expected xor, got no xor")
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf")
		if info.IsXorEnabled {
			t.Error("expected no xor, got xor")
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf.cadical_c100000.cnf.march_n810.cubes.cube8.cnf")
		cubeInfo, exists := info.Cubing.Get()
		if !exists {
			t.Error("expected cube info, got none")
		}

		if cubeInfo.Threshold != 810 {
			t.Errorf("expected cube threshold = 810, got %d", cubeInfo.Threshold)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf.cadical_c100000.cnf.march_n810.cubes.cube8.cnf")
		cubeIndex, exists := info.CubeIndex.Get()
		if !exists {
			t.Error("expected cube index, got none")
		}

		if cubeIndex != 8 {
			t.Errorf("expected cube index = 8, got %d", cubeIndex)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_ffffffffffffffffffffffffffffffff_dobbertin32.cnf.cadical_c100000.cnf.march_n810.cubes.cube8.cnf")
		simplificationInfo, exists := info.Simplification.Get()
		if !exists {
			t.Error("expected simplification info, got none")
		}

		if simplificationInfo.Simplifier != "cadical" {
			t.Errorf("expected cadical, got %s", simplificationInfo.Simplifier)
		}

		if simplificationInfo.Conflicts != 100000 {
			t.Errorf("expected cadical simplification conflicts = 100000, got %d", simplificationInfo.Conflicts)
		}
	}

	{
		info, _ := processInstanceName("nejati_encoder_md4_35_espresso_00000000000000000000000000000000_dobbertin32.cnf.satelite.cnf.march_n810.cubes.cube8.cnf")
		if info.TargetHash != "00000000000000000000000000000000" {
			t.Errorf("expected target hash = 00000000000000000000000000000000, got %s", info.TargetHash)
		}
	}
}
