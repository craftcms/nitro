package suggest

// NumberOfCPUs will suggest the number of CPUs a
// machine should be created with based on the
// current systems CPU count at runtime.
func NumberOfCPUs(num int) string {
	switch num {
	case 4:
		return "2"
	case 2:
		return "1"
	default:
		return "4"
	}
}
