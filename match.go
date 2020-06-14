package prompt

// match is a naive implementation of fuzzy search.
func match(str, pat string) bool {
	var i, j int
	for ; i < len(str) && j < len(pat); i++ {
		c := str[i]
		p := pat[j]
		if c == p {
			j++
		}
	}
	return j == len(pat)
}
