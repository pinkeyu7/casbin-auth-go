package helper

func StringContains(strSlice []string, matchStr string) bool {
	for _, each := range strSlice {
		if each == matchStr {
			return true
		}
	}
	return false
}
