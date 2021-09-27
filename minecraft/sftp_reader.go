package minecraft

import (
	"io"
	"os"
	"time"

	"github.com/hpcloud/tail"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func tailSftpFile(host, username, password, path string) (*tail.Tail, error) {

	// open an SFTP session over an existing ssh connection.
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, err
	}
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp("/tmp", "shulker-sftp-copy-")
	if err != nil {
		return nil, err
	}
	go func() {
		var offset int64 = 0
		var fileSize int64 = 0
		var first = true
		for {
			time.Sleep(1 * time.Second)
			fileInfo, err := client.Lstat(path)
			if err != nil {
				println(err.Error())
				continue
			}
			if first {
				offset = fileInfo.Size() // Set to the end to avoid sending everything that is already present
				first = false
			}

			// File replaced by minecraft server, reset the last read size
			if fileInfo.Size() < fileSize {
				offset = 0
			}

			file, err := client.Open(path)
			file.Seek(offset, io.SeekStart)
			if err == nil {
				written, _ := io.Copy(tmpFile, file)
				offset += written

				file.Close()
			}

		}
	}()

	return tailLocalFile(tmpFile.Name())
}
