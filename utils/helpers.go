package utils

func StrPtrOrDefault(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func FloatPtrOrDefault(ptr *float64) float64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}