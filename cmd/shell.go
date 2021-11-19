package cmd

import (
	"bufio"
	"os"
	"os/exec"
)

func Exec(command string, params []string) *exec.Cmd {
	cmd := exec.Command(command, params...)
	cmd.Stderr = os.Stderr
	return cmd
}

func ReadLog(filePath string, lineNumber int) ([]string, int) {
	file, _ := os.Open(filePath)
	fileScanner := bufio.NewScanner(file)
	lineCount := 1
	var lines []string
	for fileScanner.Scan() {
		if lineCount >= lineNumber {
			lines = append(lines, fileScanner.Text())
		}
		lineCount++
	}
	defer file.Close()
	return lines, lineCount - 1
}
