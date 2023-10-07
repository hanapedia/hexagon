package datageneration

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

type DummyData struct {
	ID      int    `json:"id"`
	Payload string `json:"payload"`
}

// count is the number of entries to generate
// size is the size of each data for each entry
func generateDummyData(count int, size int) []DummyData {
	data := make([]DummyData, count)
	for i := 0; i < count; i++ {
		payload, err := generateRandomString(size)
		if err != nil {
			log.Fatalf("Error generating random string %s", err)
		}
		data[i] = DummyData{
			ID:      i + 1,
			Payload: payload,
		}
	}
	return data
}

// generate random string of given size in kb
func generateRandomString(kbSize int) (string, error) {
	byteSize := kbSize * 1024
	rawByteSize := byteSize * 3 / 4
	bytes := make([]byte, rawByteSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	encoded := make([]byte, byteSize)
	base64.StdEncoding.Encode(encoded, bytes)
	return string(encoded[:byteSize]), nil
}

func generateDummyDataFile(name string, count int, size int) error {
	data := generateDummyData(count, size)

	file, err := os.Create(fmt.Sprintf("%s.json", name))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %s.json successfully", name)

	return nil
}

func main() {
	var err error
	err = generateDummyDataFile("small", 100, 1)
	if err != nil {
		log.Fatalf("Error create small dummy data %s", err)
	}
	err = generateDummyDataFile("medium", 100, 4)
	if err != nil {
		log.Fatalf("Error create medium dummy data %s", err)
	}
	err = generateDummyDataFile("large", 100, 16)
	if err != nil {
		log.Fatalf("Error create large dummy data %s", err)
	}
}
