// Copyright 2021 starship studio.
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package file control all File.
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
