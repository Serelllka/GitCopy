package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	numberOfArgs = 3
)

var destPath string
var gitBinPath string

func init() {
	var err error
	if gitBinPath, err = exec.LookPath("git"); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 3 {
		log.Fatal(fmt.Sprintf(
			"invalid number of arguments expected: %d received: %d",
			numberOfArgs,
			len(os.Args),
		))
	}

	destPath = os.Args[2]
}

func createFolderIfNotExists(path string) error {
	if f, err := exists(path); !f {
		if err != nil {
			return err
		}
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func dfs(cmd exec.Cmd, path string, hash string) {
	cmd.Args[3] = hash
	fmt.Println("command: ", cmd.String())

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	lines := strings.Split(out.String(), "\n")
	for _, item := range lines {
		info := strings.Split(item, " ")

		if len(info) < 2 {
			continue
		}
		itemType := info[1]
		fmt.Println("type: ", itemType)

		hashWithName := strings.Split(info[2], "\t")
		itemHash := hashWithName[0][:40]

		fmt.Println("hash: ", itemHash)
		itemName := ""

		if len(hashWithName) > 1 {
			itemName = hashWithName[1]
			fmt.Println("name: ", itemName)
		}
		if itemType != "tree" && itemType != "blob" {
			continue
		}

		fullPath := filepath.Join(destPath, path, itemName)

		fmt.Println("fullPath: ", fullPath)

		if itemType == "blob" {
			f, _ := os.Create(fullPath)
			defer f.Close()
			innerCmd := exec.Cmd{
				Path:   gitBinPath,
				Args:   []string{gitBinPath, "cat-file", "-p", itemHash},
				Stdout: f,
			}
			innerCmd.Run()
		}

		if itemType == "tree" {
			if f, _ := exists(fullPath); !f {
				if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
					log.Fatal(err)
				}
			}
			dfs(exec.Cmd{
				Path: cmd.Path,
				Args: cmd.Args,
			}, filepath.Join(path, itemName), itemHash)
		}
	}
}

func main() {
	hash := os.Args[1]
	err := createFolderIfNotExists(destPath)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer

	cmd := exec.Cmd{
		Path:   gitBinPath,
		Args:   []string{gitBinPath, "cat-file", "-p", hash},
		Stdout: &out,
	}

	err = cmd.Run()
	if err != nil {
		return
	}

	treeHash := strings.Split(strings.Split(out.String(), " ")[1], "\n")[0]

	dfs(exec.Cmd{
		Path: gitBinPath,
		Args: cmd.Args,
	}, "", treeHash)
}
