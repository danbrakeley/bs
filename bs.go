package bs

import "runtime"

// ExeName adds ".exe" to passed string if GOOS is windows
func ExeName(path string) string {
	if runtime.GOOS == "windows" {
		return path + ".exe"
	}
	return path
}
