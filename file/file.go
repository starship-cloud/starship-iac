package file

import (
	"bufio"
	"io"
	"os"
	"os/exec"
)

func WriteFileByCmd(fileName string, dir string, cmd *exec.Cmd) error {
	f, err := createFile(fileName, dir)
	cmd.Dir = dir
	stdout, err := cmd.StdoutPipe()
	err = cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		print(line)
		_, err = f.WriteString(line)
		f.Sync()
	}
	err = cmd.Wait()
	return err
}

func createFile(fileName string, dir string) (*os.File, error) {
	os.MkdirAll(dir, 755)
	return os.Create(dir + fileName)
}
