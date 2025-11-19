package utils

func InList[T comparable](key T, list []T) bool {
	for _, s := range list {
		if key == s {
			return true
		}
	}
	return false
}
