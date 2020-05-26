package suggest

// NumberOfCPUs will suggest the number of CPUs a
// machine should be created with based on the
// current systems CPU count at runtime.
func NumberOfCPUs(num int) string {
	// on 8 core systems, the num from runtime.NumCPU() will return
	// 18. So we need to divide the number of CPUs by half
	switch num / 2 {
	case 4:
		return "2"
	case 2:
		return "1"
	case 1:
		return "1"
	}

	// return 4 by default
	return "4"
}
