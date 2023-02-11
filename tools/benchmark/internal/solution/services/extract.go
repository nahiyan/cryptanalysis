package services

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func (solutionSvc *SolutionService) ExtractMessage(solutionLiterals []int) ([]byte, error) {
	var messageBuilder strings.Builder
	for i := 0; i < 16; i++ {
		word := solutionLiterals[i*32 : i*32+32]
		word_ := lo.Reduce(word, func(acc string, bit int, _ int) string {
			if bit > 1 {
				return "1" + acc
			} else {
				return "0" + acc
			}
		}, "")
		word__, err := strconv.ParseUint(word_, 2, 32)
		if err != nil {
			return nil, err
		}
		messageBuilder.WriteString(fmt.Sprintf("%08x", word__))
	}

	message := messageBuilder.String()
	bytes, err := hex.DecodeString(message)
	if err != nil {
		return nil, err
	}
	return bytes, err
}
