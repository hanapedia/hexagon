package utils

func Int32Ptr(i int32) *int32 { return &i }

func Btof64(b bool) float64 {
	if b {
		return float64(1)
	}
	return float64(0)
}

func Btos(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func GetMapKeys[v interface{}](m map[string]v) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
