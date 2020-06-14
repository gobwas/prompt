// +build !aix,!darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package prompt

func prefillInput(s string) error {
	// TODO: support this on other platforms.
	return nil
}
