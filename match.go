package prompt

// match is a naive implementation of fuzzy search.
func match(str, pat string) bool {
	const toLower = 'a' - 'A'
	var i, j int
	for ; i < len(str) && j < len(pat); i++ {
		c := str[i] | toLower
		p := pat[j] | toLower
		if c == p {
			j++
		}
	}
	return j == len(pat)
}
