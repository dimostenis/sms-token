package main

import (
	"flag"
	"smstoken/pkg/symlinks"
	"smstoken/pkg/token"
)

func main() {
	// parse args
	var install bool
	var uninstall bool
	flag.BoolVar(&install, "install", false, "create symlink in /usr/local/bin so 'token' is system-wide")
	flag.BoolVar(&uninstall, "uninstall", false, "delete existing symlink (if exists) in /usr/local/bin")
	flag.Parse()

	if install {
		// user wants to create symlink
		symlinks.DeleteSymlink()
		symlinks.CreateSymlink()
	} else if uninstall {
		// user wants just to delete existing symlink
		symlinks.DeleteSymlink()
	} else {
		// no arguments, default behaviour
		token.GetToken()
	}
}
