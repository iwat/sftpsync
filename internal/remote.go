package internal

import (
	"fmt"
	"path/filepath"

	"github.com/kr/fs"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/iwat/go-log"
)

func NewSFTPClient(addr string, config *ssh.ClientConfig) (*sftp.Client, error) {
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return sftp.NewClient(conn)
}

func BuildClients(host string, config *ssh.ClientConfig, num int) []*sftp.Client {
	clients := make([]*sftp.Client, 4)
	for i := 0; i < 4; i++ {
		client, err := NewSFTPClient(host, config)
		if err != nil {
			log.ERR.Fatalln("could not connect sftp:", err)
		}

		clients[i] = client
	}

	return clients
}

func BuildRemoteFileList(walker *fs.Walker, basepath string) map[string]file {
	output := make(map[string]file)

	i := 0

	for walker.Step() {
		if err := walker.Err(); err != nil {
			log.WRN.Println("walker error:", err)
			continue
		}

		fmt.Print(".")
		i++

		if i > 80 {
			i = 0
			fmt.Println()
		}

		rel, err := filepath.Rel(basepath, walker.Path())
		if err != nil {
			log.WRN.Println("could not resolve relative path:", walker.Path())
			continue
		}

		if rel == "." {
			continue
		}

		stat := walker.Stat()
		file := file{
			mode:    stat.Mode(),
			size:    stat.Size(),
			mod:     stat.ModTime(),
			path:    walker.Path(),
			relPath: rel,
		}

		output[file.relPath] = file
	}

	if i > 0 {
		fmt.Println()
	}

	return output
}
