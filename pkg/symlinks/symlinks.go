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

var repoPath string = getRepoPath()

type Arch string

const (
	ARM64 Arch = "arm64"
	AMD64 Arch = "amd64"
)

func getRepoPath() string {
	p, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return p
}

func getScriptContent() string {
	shebang := "#!/bin/zsh"
	scriptDir := repoPath + "/bin"
	scriptName := fmt.Sprintf("/token-%s-darwin", getArch())
	script := path.Join(scriptDir, scriptName)
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

func writeFile(scriptPath string, content string) {
	// prep parent directory, which may not be pushed in repo
	dir, _ := path.Split(scriptPath)
	mkdir(dir)

	// write file
	err := os.WriteFile(scriptPath, []byte(content), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createSymlink(sourceFile string, targetFile string) {
	// create symbolic link to file (its not executable)
	err := os.Symlink(sourceFile, targetFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func makeFileExecutable(target_file string) {
	// make the symbolic link executable
	err := os.Chmod(target_file, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getArch() Arch {
	if runtime.GOARCH == "arm64" {
		return ARM64
	} else {
		return AMD64 // default arch, shall work everywhere
	}
}

func deleteFile(path string) {
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
	scriptPath := path.Join(repoPath, "bin/token")

	// get shell script content
	scriptContent := getScriptContent()

	// write script
	writeFile(scriptPath, scriptContent)

	// create symlink so its system-wide
	createSymlink(scriptPath, symlink)
	makeFileExecutable(symlink)

	fmt.Printf("Symlink '%s' created\n", symlink)
}

func DeleteSymlink() {
	deleteFile(symlink)

	fmt.Printf("Symlink '%s' deleted\n", symlink)
}
