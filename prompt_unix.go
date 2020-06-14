// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package prompt

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func prefillInput(s string) error {
	for i := 0; i < len(s); i++ {
		c := s[i]
		_, _, eno := unix.RawSyscall(
			unix.SYS_IOCTL,
			os.Stdin.Fd(),
			unix.TIOCSTI,
			uintptr(unsafe.Pointer(&c)))
		if eno != 0 {
			return fmt.Errorf("%v", eno)
		}
	}
	return nil
}
