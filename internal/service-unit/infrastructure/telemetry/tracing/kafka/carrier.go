package kafka

import "github.com/segmentio/kafka-go"

type KafkaCarrier struct {
    Headers []kafka.Header
}

func (c *KafkaCarrier) Get(key string) string {
	for _, h := range c.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c *KafkaCarrier) Set(key string, value string) {
	for i, h := range c.Headers {
		if h.Key == key {
			c.Headers[i].Value = []byte(value)
			return
		}
	}

	// If the key doesn't exist, add a new header.
	c.Headers = append(c.Headers, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (c *KafkaCarrier) Keys() []string {
	keys := make([]string, 0, len(c.Headers))
	for _, h := range c.Headers {
		keys = append(keys, h.Key)
	}
	return keys
}
