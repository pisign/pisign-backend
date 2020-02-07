package api

func list() []string {
	// TODO: utilize specfile
	return []string{"weather"}
}

// Valid returns a list of all current valid go api modules
func Valid(name string) bool {
	// TODO: Make this more efficient
	for _, n := range list() {
		if n == name {
			return true
		}
	}
	return false
}
