package tool

func SliceUnique(slice []string) (uniqueslice []string) {
	for _, v := range slice {
		if !InSliceIface(v, uniqueslice) {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}

// InSliceIface checks given interface in interface slice.
func InSliceIface(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
