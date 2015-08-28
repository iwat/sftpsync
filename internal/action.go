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

		err := client.Remove(d.path)
		if err != nil {
			log.WRN.Print("DEL ", d, ": ", err)
			continue
		}

		log.NFO.Println("DEL", d)
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

	err := client.Mkdir(basepath + "/" + p.relPath)
	if err != nil {
		log.WRN.Print("MKDIR ", p, ": ", err)
		return
	}

	log.NFO.Print("MKDIR ", p)
}

func putFile(p file, client *sftp.Client, basepath string, dryRun bool) {
	if dryRun {
		log.NFO.Print("!PUT ", p)
		return
	}

	remote, err := client.Create(basepath + "/" + p.relPath)
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

	_, err = io.Copy(remote, local)
	if err != nil {
		log.WRN.Print("PUT ", p, ": ", err)
		return
	}

	log.NFO.Println("PUT", p)
}
