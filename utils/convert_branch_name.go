package utils

func ConvertBranchName(branchName string) string {
	branchNameBytes := []byte(branchName)
	i := 0
	for _, b := range branchNameBytes {
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') || b == '-' {
			branchNameBytes[i] = b
			i++
		} else {
			branchNameBytes[i] = '-'
			i++
		}
	}

	return string(branchNameBytes[:i])
}
