package utils

import (
	mrand "math/rand"
	"strings"
)

const LoremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla ultrices mollis finibus. Nullam justo tellus, blandit quis urna quis, aliquam luctus dui. Morbi luctus dolor magna, a dictum tortor elementum eget. Vestibulum finibus faucibus commodo. Suspendisse potenti. Cras vulputate ultrices metus at dignissim. Suspendisse porta lectus ipsum, quis faucibus ligula venenatis in. Nam et nisl tellus. In vehicula vitae orci ut dignissim."

func GenerateRandomString(byteSize int64) string {
    // Create a string builder for efficient string concatenation
    var sb strings.Builder

    // Iterate and append characters until we reach the desired byte size
    for i := 0; i < int(byteSize); i++ {
        sb.WriteByte(LoremIpsum[i%len(LoremIpsum)])
    }

    return sb.String()
}

// generates random integer from 1 to 100
func GetRandomId(min int, max int) int {
	if min >= max {
		max = 100
		min = 1
	}

	randomInt := mrand.Intn(max) + min
	return randomInt
}
