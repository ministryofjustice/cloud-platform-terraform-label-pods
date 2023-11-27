package utils

func Contains(sysNs []string, ns string) bool {
	for _, sys := range sysNs {
		if sys == ns {
			return true
		}
	}
	return false
}
