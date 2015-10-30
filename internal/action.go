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
			log.NFO.Println("!DEL", d)
			continue
		}

		log.NFO.Println("DEL", d)
		err := client.Remove(d.path)
		if err != nil {
			log.WRN.Print("DEL ", d, ": ", err)
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
		log.NFO.Print("!MKDIR ", p)
		return
	}

	log.NFO.Print("MKDIR ", p)
	err := client.Mkdir(basepath + "/" + p.relPath)
	if err != nil {
		log.WRN.Print("MKDIR ", p, ": ", err)
		return
	}
}

func putFile(p file, client *sftp.Client, basepath string, dryRun bool) {
	if dryRun {
		if p.offset > 0 {
			log.NFO.Print("!APPEND", p)
		} else {
			log.NFO.Print("!PUT ", p)
		}
		return
	}

	var remote *sftp.File
	var err error
	if p.offset > 0 {
		log.NFO.Println("APPEND", p)
		remote, err = client.OpenFile(basepath+"/"+p.relPath, os.O_APPEND)
	} else {
		log.NFO.Println("PUT", p)
		remote, err = client.Create(basepath + "/" + p.relPath)
	}
	if err != nil {
		log.WRN.Print("PUT ", p, ": ", err)
		return
	}

	defer remote.Close()

	local, err := os.Open(p.path)
	if err != nil {
		log.WRN.Print("PUT ", p, ": ", err)
		return
	}

	defer local.Close()

	if p.offset > 0 {
		local.Seek(p.offset, os.SEEK_SET)
	}

	_, err = io.Copy(remote, local)
	if err != nil {
		log.WRN.Print("PUT ", p, ": ", err)
		return
	}
}
