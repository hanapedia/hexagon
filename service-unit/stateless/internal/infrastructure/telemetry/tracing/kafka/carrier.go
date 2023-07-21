package kafka

import "github.com/segmentio/kafka-go"

type KafkaCarrier []kafka.Header

func (c KafkaCarrier) Get(key string) string {
	for _, h := range c {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c KafkaCarrier) Set(key string, value string) {
	for i, h := range c {
		if h.Key == key {
			c[i].Value = []byte(value)
			return
		}
	}

	// If the key doesn't exist, add a new header.
	c = append(c, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (c KafkaCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for _, h := range c {
		keys = append(keys, h.Key)
	}
	return keys
}
