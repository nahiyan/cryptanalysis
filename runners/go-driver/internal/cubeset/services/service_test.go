package services

import (
	"encoding/binary"
	"os"
	"testing"
)

func readAt(bCubesFile *os.File, index int64, t *testing.T) int32 {
	bytes := make([]byte, 4)
	_, err := bCubesFile.ReadAt(bytes, index*4)
	if err != nil {
		t.Fatal(err)
	}
	return int32(binary.BigEndian.Uint32(bytes))
}

func initialize(t *testing.T) *CubesetService {
	svc := CubesetService{}

	cubeset := "a -414 -408 -250 232 -2234 2487 2248 -421 -494 -373 0\na -414 -408 -250 232 -2234 2487 -2248 -359 2238 226 587 595 531 532 2290 373 2258 -550 0\na 414 -415 37358 250 -40233 232 473 234 2487 545 249 -367 469 238 421 -2248 490 492 -256 375 -377 -498 0\n"
	err := os.WriteFile("x.cubes", []byte(cubeset), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = svc.BinEncode("x.cubes")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("x.cubes.bcubes")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("x.cubes.bcubes.map")
	if err != nil {
		t.Fatal(err)
	}

	return &svc
}

func cleanup() {
	os.Remove("x.cubes")
	os.Remove("x.cubes.bcubes")
	os.Remove("x.cubes.bcubes.map")
}

func TestBinEncode(t *testing.T) {
	initialize(t)
	defer cleanup()

	bCubesFile, err := os.OpenFile("x.cubes.bcubes", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	{
		literal := readAt(bCubesFile, 0, t)
		if literal != -414 {
			t.Fatalf("expected -414 but got %d", literal)
		}
	}
	{
		literal := readAt(bCubesFile, 4, t)
		if literal != -2234 {
			t.Fatalf("expected -2234 but got %d", literal)
		}
	}
	{
		literal := readAt(bCubesFile, 9, t)
		if literal != -373 {
			t.Fatalf("expected -373 but got %d", literal)
		}
	}
	{
		literal := readAt(bCubesFile, 10, t)
		if literal != -414 {
			t.Fatalf("expected -414 but got %d", literal)
		}
	}
	{
		literal := readAt(bCubesFile, 13, t)
		if literal != 232 {
			t.Fatalf("expected 232 but got %d", literal)
		}
	}
	{
		literal := readAt(bCubesFile, 27, t)
		if literal != -550 {
			t.Fatalf("expected -550 but got %d", literal)
		}
	}
	{
		literal := readAt(bCubesFile, 49, t)
		if literal != -498 {
			t.Fatalf("expected -498 but got %d", literal)
		}
	}
	{
		b := make([]byte, 4)
		_, err := bCubesFile.ReadAt(b, 50*4)
		if err == nil || err.Error() != "EOF" {
			t.Fatalf("expected EOF but got %s", err)
		}
	}

	mapFile, err := os.OpenFile("x.cubes.bcubes.map", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}

	{
		bytes := make([]byte, 4)
		_, err = mapFile.ReadAt(bytes, 0)
		if err != nil {
			t.Fatal(err)
		}
		index := binary.BigEndian.Uint32(bytes)
		if index != 10 {
			t.Fatalf("expected 10 but got %d", index)
		}
	}
	{
		bytes := make([]byte, 4)
		_, err = mapFile.ReadAt(bytes, 4)
		if err != nil {
			t.Fatal(err)
		}
		index := binary.BigEndian.Uint32(bytes)
		if index != 28 {
			t.Fatalf("expected 28 but got %d", index)
		}
	}
	{
		bytes := make([]byte, 4)
		_, err = mapFile.ReadAt(bytes, 8)
		if err != nil {
			t.Fatal(err)
		}
		index := binary.BigEndian.Uint32(bytes)
		if index != 50 {
			t.Fatalf("expected 50 but got %d", index)
		}
	}
}

func areCubesEqual(cube []int, expectedCube []int) bool {
	for i, literal := range cube {
		if literal != expectedCube[i] {
			return false
		}
	}

	return true
}

func TestGetCube(t *testing.T) {
	svc := initialize(t)
	defer cleanup()

	{
		cube, err := svc.GetCube("x.cubes", 1)
		if err != nil {
			t.Fatal(err)
		}
		expectedCube := []int{-414, -408, -250, 232, -2234, 2487, 2248, -421, -494, -373}
		if !areCubesEqual(cube, expectedCube) {
			t.Fatalf("got %v but expected %v", cube, expectedCube)
		}
	}

	{
		cube, err := svc.GetCube("x.cubes", 2)
		if err != nil {
			t.Fatal(err)
		}
		expectedCube := []int{-414, -408, -250, 232, -2234, 2487, -2248, -359, 2238, 226, 587, 595, 531, 532, 2290, 373, 2258, -550}
		if !areCubesEqual(cube, expectedCube) {
			t.Fatalf("got %v but expected %v", cube, expectedCube)
		}
	}
	{
		cube, err := svc.GetCube("x.cubes", 3)
		if err != nil {
			t.Fatal(err)
		}
		expectedCube := []int{414, -415, 37358, 250, -40233, 232, 473, 234, 2487, 545, 249, -367, 469, 238, 421, -2248, 490, 492, -256, 375, -377, -498}
		if !areCubesEqual(cube, expectedCube) {
			t.Fatalf("got %v but expected %v", cube, expectedCube)
		}
	}
}
