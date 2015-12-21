// sftpsync - Sync local file system to SFTP

// Copyright (c) 2015 Chaiwat Shuetrakoonpaiboon. All rights reserved.
//
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

package internal

import (
	"fmt"
	"path/filepath"

	"github.com/kr/fs"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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
			log.Crit("could not connect sftp", "err", err)
			return nil
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
			log.Warn("walker error", "err", err)
			continue
		}

		fmt.Print(".")
		i++

		if i >= 80 {
			i = 0
			fmt.Println()
		}

		rel, err := filepath.Rel(basepath, walker.Path())
		if err != nil {
			log.Warn("could not resolve relative path", "path", walker.Path())
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
