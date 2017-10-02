/*
	Package entrypoint replaces the shell version of
		exec $@
	which is often use when creating docker entrypoint scripts to wrap a
	binary with some initial setup.
*/
package entrypoint

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func Exec() {
	flag.Parse()
	if len(os.Args) == 1 {
		return
	}
	cmd, err := exec.LookPath(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := syscall.Exec(cmd, flag.Args(), os.Environ()); err != nil {
		log.Fatal(err)
	}
}
