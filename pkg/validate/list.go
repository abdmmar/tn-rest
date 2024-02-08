package validate

func Include[K comparable](list []K) bool {
	seen := make(map[K]bool)

	for _, str := range list {
		if _, has := seen[str]; has {
			return false
		}
		seen[str] = true
	}

	return true
}
