package symlinks

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"syscall"
)

// this is default path to system-wide executable
const symlink string = "/usr/local/bin/token"

var repo_path string = get_repo_path()

type Arch string

const (
	ARM64 Arch = "arm64"
	AMD64 Arch = "amd64"
)

func get_repo_path() string {
	p, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return p
}

func get_script_content() string {
	shebang := "#!/bin/zsh"
	script_dir := repo_path + "/bin"
	script_name := fmt.Sprintf("/token-%s-darwin", get_arch())
	script := path.Join(script_dir, script_name)
	qry := "sqlite3 -line ~/Library/Messages/chat.db \"SELECT m.text, m.date FROM message m WHERE text LIKE '%code%' ORDER BY m.ROWID DESC LIMIT 1\""
	content := fmt.Sprintf("%s\n%s | %s\n", shebang, qry, script)
	return content
}

func mkdir(path string) {
	err := os.Mkdir(path, 0644)
	if err != nil {
		if errors.Is(err, syscall.EEXIST) {
			// nothing to do
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}

func write_file(script_path string, content string) {
	// prep parent directory, which may not be pushed in repo
	dir, _ := path.Split(script_path)
	mkdir(dir)

	// write file
	err := os.WriteFile(script_path, []byte(content), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func create_symlink(source_file string, target_file string) {
	// create symbolic link to file (its not executable)
	err := os.Symlink(source_file, target_file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func make_file_executable(target_file string) {
	// make the symbolic link executable
	err := os.Chmod(target_file, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func get_arch() Arch {
	if runtime.GOARCH == "arm64" {
		return ARM64
	} else {
		return AMD64 // default arch, shall work everywhere
	}
}

func delete_file(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// nothing to do
	} else if err != nil {
		fmt.Println(err)
	} else {
		os.Remove(path)
	}
}

func CreateSymlink() {
	// find script path
	script_path := path.Join(repo_path, "bin/token")

	// get shell script content
	script_content := get_script_content()

	// write script
	write_file(script_path, script_content)

	// create symlink so its system-wide
	create_symlink(script_path, symlink)
	make_file_executable(symlink)

	fmt.Printf("Symlink '%s' created\n", symlink)
}

func DeleteSymlink() {
	delete_file(symlink)

	fmt.Printf("Symlink '%s' deleted\n", symlink)
}
