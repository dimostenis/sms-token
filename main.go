package main

import (
	"flag"
	"fmt"
	"os"
	"smstoken/pkg/sms"
	"smstoken/pkg/symlinks"
)

func main() {
	// parse args
	var install bool
	var uninstall bool
	var text string
	flag.StringVar(&text, "text", "", "fake SMS text to test SMS parsing")
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
		const tenMinutes = 10 * 60
		config := sms.Config{
			MaxAge: tenMinutes,
		}
		if text != "" {
			config.Debug = true
			config.DebugText = text
			config.MaxAge = 100000 * 100000 // 300+ years
		}
		_, err := sms.ReadToken(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
