package easytmpl

// IsBlank checks if a byte slice is blank (contains only spaces).
// It returns true if the byte slice is blank, otherwise false.
// An empty byte slice is considered blank.
// For example:
//
//	IsBlank([]byte("   ")) // returns true
//	IsBlank(nil)    // returns true
//	IsBlank([]byte(" a ")) // returns false
func IsBlank(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] != byte(' ') {
			return false
		}
	}
	return true
}
