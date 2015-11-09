package internal

import (
	"io"
	"os"

	"github.com/pkg/sftp"

	"github.com/iwat/go-log"
)

func ProcessDelete(client *sftp.Client, dc <-chan file, dryRun bool, done chan<- bool) {
	for d := range dc {
		if dryRun {
			log.Infoln("!DEL", d)
			continue
		}

		log.Infoln("DEL", d)
		err := client.Remove(d.path)
		if err != nil {
			log.Warn("DEL ", d, ": ", err)
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
		log.Info("!MKDIR ", p)
		return
	}

	log.Info("MKDIR ", p)
	err := client.Mkdir(basepath + "/" + p.relPath)
	if err != nil {
		log.Warn("MKDIR ", p, ": ", err)
		return
	}
}

func putFile(p file, client *sftp.Client, basepath string, dryRun bool) {
	if dryRun {
		if p.offset > 0 {
			log.Info("!APPEND", p)
		} else {
			log.Info("!PUT ", p)
		}
		return
	}

	var remote *sftp.File
	var err error
	if p.offset > 0 {
		log.Infoln("APPEND", p)
		remote, err = client.OpenFile(basepath+"/"+p.relPath, os.O_RDWR|os.O_APPEND)
	} else {
		log.Infoln("PUT", p)
		remote, err = client.Create(basepath + "/" + p.relPath)
	}
	if err != nil {
		log.Warn("PUT ", p, ": ", err)
		return
	}

	defer remote.Close()

	local, err := os.Open(p.path)
	if err != nil {
		log.Warn("PUT ", p, ": ", err)
		return
	}

	defer local.Close()

	if p.offset > 0 {
		_, err := local.Seek(p.offset, os.SEEK_SET)
		if err != nil {
			log.Warnln("Seek error:", err)
			return
		}
	}

	_, err = io.Copy(remote, local)
	if err != nil {
		log.Warn("PUT ", p, ": ", err)
		return
	}
}
