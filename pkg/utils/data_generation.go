package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strconv"

	"github.com/samber/lo"
)

const (
	digits  = "0123456789"
	special = "!@#$%^&*()-_=+[]{}|;:',.<>?/"
	upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower   = "abcdefghijklmnopqrstuvwxyz"
)

func GenerateUpperCaseString(length int) string {
	charset := []rune(upper)

	return lo.RandomString(length, charset)
}

func GeneratePassword(length int) string {
	charset := []rune(digits + special + upper + lower)

	requiredChars := []byte{
		digits[GenerateInt64(0, int64(len(digits))-1)],
		special[GenerateInt64(0, int64(len(special))-1)],
		upper[GenerateInt64(0, int64(len(upper))-1)],
		lower[GenerateInt64(0, int64(len(lower))-1)],
	}
	remainingChars := lo.RandomString(length-len(requiredChars), charset)

	return string(append(requiredChars, []byte(remainingChars)...))
}

func GenerateInt64(minVal, maxVal int64) int64 {
	if minVal > maxVal {
		panic("minVal must be less or equal to maxVal")
	}

	rangeSize := maxVal - minVal + 1

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		randomNumber = big.NewInt(mrand.Int63n(rangeSize)) // nolint
	}

	return minVal + randomNumber.Int64()
}

func GenerateFloat64(minVal, maxVal float64) float64 {
	if minVal > maxVal {
		panic("minVal must be less or equal to maxVal")
	}

	intMin := int64(0)
	scaleFactor := int64(1 << 53)
	intMax := scaleFactor

	randomInt := GenerateInt64(intMin, intMax-1)
	randomFloat := float64(randomInt) / float64(scaleFactor)

	return minVal + randomFloat*(maxVal-minVal)
}

func GenerateUniqueNumbers(n, minVal, maxVal int64) []int64 {
	if minVal > maxVal {
		panic("minVal must be less or equal to maxVal")
	}
	if n < 0 {
		panic("n must be a positive number")
	}
	if n > maxVal-minVal+1 {
		panic(fmt.Sprintf("range is not large enough to generate %d unique numbers", n))
	}

	uniqueNumbers := make(map[int64]any)
	numbers := make([]int64, 0, n)

	for int64(len(numbers)) < n {
		num := GenerateInt64(minVal, maxVal)
		if _, exists := uniqueNumbers[num]; !exists {
			uniqueNumbers[num] = nil
			numbers = append(numbers, num)
		}
	}

	return numbers
}

func GenerateFlatMap(numPairs int) map[string]any {
	obj := make(map[string]any)

	for i := range numPairs {
		key := "key_" + strconv.Itoa(i)

		var value any
		switch i % 3 {
		case 0:
			value = "val_" + strconv.Itoa(i)
		case 1:
			value = i
		case 2:
			value = i%2 == 0
		}

		obj[key] = value
	}

	return obj
}

func GenerateFlatJSON(numPairs int) string {
	obj := GenerateFlatMap(numPairs)

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic("cannot marshal generated map")
	}
	return string(jsonBytes)
}
