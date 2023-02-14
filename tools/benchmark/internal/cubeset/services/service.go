package services

import (
	"bufio"
	"encoding/binary"
	"errors"
	"os"
	"strconv"
	"strings"
)

func (cubesetSvc *CubesetService) BinEncode(cubesetPath string) error {
	cubesFile, err := os.Open(cubesetPath)
	if err != nil {
		return err
	}
	defer cubesFile.Close()

	binCubesMapFile, err := os.OpenFile(cubesetPath+".bcubes.map", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer binCubesMapFile.Close()

	binCubesFile, err := os.OpenFile(cubesetPath+".bcubes", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer binCubesFile.Close()

	scanner := bufio.NewScanner(cubesFile)
	addressAccumulator := uint32(0)
	binCubesMapWriter := bufio.NewWriter(binCubesMapFile)
	binCubesWriter := bufio.NewWriter(binCubesFile)
	for scanner.Scan() {
		// Bin. Cube
		line := scanner.Text()
		if line == "a 0" {
			return errors.New("empty cubeset file")
		}
		literals := strings.Fields(line[2 : len(line)-2])
		for _, literal_ := range literals {
			literal, err := strconv.ParseInt(literal_, 10, 16)
			if err != nil {
				return err
			}
			bytes := make([]byte, 2)
			binary.BigEndian.PutUint16(bytes, uint16(literal))
			_, err = binCubesWriter.Write(bytes)
			if err != nil {
				return err
			}
		}

		// Line map
		{
			literalsCount := uint32(len(literals))
			address := addressAccumulator + literalsCount
			addressAccumulator = address

			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, address)
			_, err := binCubesMapWriter.Write(bytes)
			if err != nil {
				return err
			}
		}
	}
	binCubesWriter.Flush()
	binCubesMapWriter.Flush()

	return nil
}

func (cubesetSvc *CubesetService) GetCube(cubesetPath string, index int) ([]int, error) {
	cubesFile, err := os.OpenFile(cubesetPath+".bcubes", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer cubesFile.Close()

	mapFile, err := os.OpenFile(cubesetPath+".bcubes.map", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer mapFile.Close()

	// Every index in the map file is the index of the word in the bcubes file, and every word in the bcubes file is 16 bits in size
	endIndex_ := make([]byte, 4)
	_, err = mapFile.ReadAt(endIndex_, int64((index-1)*4))
	if err != nil {
		return nil, err
	}
	endIndex := binary.BigEndian.Uint32(endIndex_)
	endAddress := int64(endIndex * 2)

	startIndex := uint32(0)
	if index > 1 {
		startAddress_ := make([]byte, 4)
		_, err := mapFile.ReadAt(startAddress_, int64(index-2)*4)
		if err != nil {
			return nil, err
		}
		startIndex = binary.BigEndian.Uint32(startAddress_)
	}
	startAddress := int64(startIndex * 2)

	bytes := make([]byte, endAddress-startAddress)
	_, err = cubesFile.ReadAt(bytes, startAddress)
	if err != nil {
		return nil, err
	}

	cube := []int{}
	for i := 0; i < len(bytes); i += 2 {
		word := bytes[i : i+2]
		literal := int16(binary.BigEndian.Uint16(word))
		cube = append(cube, int(literal))
	}

	return cube, nil
}
