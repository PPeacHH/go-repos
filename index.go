package main

import (
	"bufio"
	//"bufio"
	//"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func readFile(path string) ([]byte, error) {
	parentPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fullPath := filepath.Join(parentPath, path)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return read(file)
}

func read(file *os.File) ([]byte, error) {
	br := bufio.NewReader(file)
	br.ReadLine()
	return _, nil
}

func main() {
	path := "index.go"
	ba, err := readFile(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("The content of '%s':\n%s\n", path, ba)
}