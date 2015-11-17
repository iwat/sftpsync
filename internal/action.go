package internal

import (
	"io"
	"os"

	"github.com/pkg/sftp"
)

func ProcessDelete(client *sftp.Client, dc <-chan file, dryRun bool, done chan<- bool) {
	for d := range dc {
		if dryRun {
			log.Info("!DEL", "path", d.relPath, "size", d.size)
			continue
		}

		log.Info("DEL", "path", d.relPath)
		err := client.Remove(d.path)
		if err != nil {
			log.Warn("DEL", "path", d.relPath, "err", err)
			continue
		}
	}

	done <- true
}

func ProcessPut(client *sftp.Client, basepath string, pc <-chan file, dryRun bool, done chan<- bool) {
	for p := range pc {
		if p.mode.IsDir() {
			putDir(p, client, basepath, dryRun)
		} else {
			putFile(p, client, basepath, dryRun)
		}
	}

	done <- true
}

func putDir(p file, client *sftp.Client, basepath string, dryRun bool) {
	if dryRun {
		log.Info("!MKDIR", "path", p.relPath)
		return
	}

	log.Info("MKDIR", "path", p.relPath)
	err := client.Mkdir(basepath + "/" + p.relPath)
	if err != nil {
		log.Warn("MKDIR", "path", p.relPath, "err", err)
		return
	}
}

func putFile(p file, client *sftp.Client, basepath string, dryRun bool) {
	if dryRun {
		if p.offset > 0 {
			log.Info("!APPEND", "path", p.relPath, "size", p.size)
		} else {
			log.Info("!PUT", "path", p.relPath, "size", p.size)
		}
		return
	}

	var remote *sftp.File
	var err error
	if p.offset > 0 {
		log.Info("APPEND", "path", p.relPath, "size", p.size)
		remote, err = client.OpenFile(basepath+"/"+p.relPath, os.O_RDWR|os.O_APPEND)
	} else {
		log.Info("PUT", "path", p.relPath, "size", p.size)
		remote, err = client.Create(basepath + "/" + p.relPath)
	}
	if err != nil {
		log.Warn("PUT", "path", p.relPath, "size", p.size, "err", err)
		return
	}

	defer remote.Close()

	local, err := os.Open(p.path)
	if err != nil {
		log.Warn("PUT", "path", p.relPath, "size", p.size, "err", err)
		return
	}

	defer local.Close()

	if p.offset > 0 {
		_, err := local.Seek(p.offset, os.SEEK_SET)
		if err != nil {
			log.Warn("Seek error", "err", err)
			return
		}
	}

	_, err = io.Copy(remote, local)
	if err != nil {
		log.Warn("PUT", "path", p.relPath, "size", p.size, "err", err)
		return
	}
}
