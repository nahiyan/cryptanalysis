package services

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func (solutionSvc *SolutionService) ExtractFromLiterals(literals []int) ([]byte, error) {
	var messageBuilder strings.Builder
	for i := 0; i < len(literals)/32; i++ {
		word := literals[i*32 : i*32+32]
		word_ := lo.Reduce(word, func(acc string, bit int, _ int) string {
			if bit > 0 {
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
