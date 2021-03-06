// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package posix provides a few POSIX-based facilities for local and remote scripting
package posix

import (
	"io/ioutil"
	"os/exec"
)

// Exec executes prog with working directory dir, flags argv, and standard input stdin.
// Exec returns the standard output and error streams as strings.
func Exec(prog, dir string, stdin string, argv ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(prog, argv...)
	cmd.Dir = dir

	stdinWriter, err := cmd.StdinPipe()
	if err != nil {
		return "", "", err
	}
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}
	if err := cmd.Start(); err != nil {
		return "", "", err
	}
	// Since Run is meant for non-interactive execution, we pump all the stdin first,
	// then we (sequentially) read all of stdout and then all of stderr.
	// XXX: Is it possible to block if the program's stderr buffer fills while we are
	// consuming the stdout?
	_, err = stdinWriter.Write([]byte(stdin))
	if err != nil {
		return "", "", err
	}
	err = stdinWriter.Close()
	if err != nil {
		return "", "", err
	}
	stdoutBuf, _ := ioutil.ReadAll(stdoutReader)
	stderrBuf, _ := ioutil.ReadAll(stderrReader)

	return string(stdoutBuf), string(stderrBuf), cmd.Wait()
}

// Shell executes sh and feeds it standard input shellStdin.
func Shell(shellStdin string) (stdout, stderr string, err error) {
	return Exec("sh", "", shellStdin)
}

// RemoteShell executes sh on remoteHost via ssh and feeds it remoteShellStdin on the standard input.
func RemoteShell(remoteHost, remoteShellStdin string) (stdout, stderr string, err error) {
	return Exec("ssh", "", remoteShellStdin, remoteHost, "sh -il")
}

// DownloadDir copies the contents of remoteDir on remoteHost to local directory sourceDir, using rsync over ssh.
func DownloadDir(remoteHost, remoteDir, sourceDir string) error {
	_, _, err := Exec("rsync", "", "", "-acrv", "--rsh=ssh", remoteHost+":"+remoteDir+"/", sourceDir+"/")
	return err
}

// UploadDir copies the contents of sourceDir recursively into remoteDir.
// remoteDir must be present on the remote host.
func UploadDir(remoteHost, sourceDir, remoteDir string) error {
	_, _, err := Exec("rsync", "", "", "-acrv", "--rsh=ssh", sourceDir+"/", remoteHost+":"+remoteDir+"/")
	return err
}
