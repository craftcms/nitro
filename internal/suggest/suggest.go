package suggest

// NumberOfCPUs will suggest the number of CPUs a
// machine should be created with based on the
// current systems CPU count at runtime
func NumberOfCPUs(num int) string {
	// Go counts the number of logical CPUs with runtime.NumCPU()
	// we are only concerned about low end CPUs where performance
	// is really impacted
	switch num {
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
