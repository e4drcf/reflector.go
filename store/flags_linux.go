// +build linux

package store

import (
	"os"
)

var openFileFlags = os.O_WRONLY | os.O_CREATE | syscall.O_DIRECT
