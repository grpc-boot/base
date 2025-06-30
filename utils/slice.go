package utils

func Unique[T comparable](arr []T) []T {
	if len(arr) < 2 {
		return arr
	}

	write := 0
	for i, _ := range arr {
		duplicate := false
		for j := 0; j < write; j++ {
			if arr[j] == arr[i] {
				duplicate = true
				break
			}
		}

		if !duplicate {
			arr[write] = arr[i]
			write++
		}
	}

	return arr[:write]
}
