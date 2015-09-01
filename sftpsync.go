package sftpsync

import (
	"net"
	"net/url"
	"strconv"

	"golang.org/x/crypto/ssh"

	"github.com/iwat/go-log"
	"github.com/iwat/go-seckeychain"
	"github.com/iwat/sftpsync/internal"
)

type SyncManager struct {
	DryRun bool
	Local  string
	Remote string

	SSHClientConfig *ssh.ClientConfig

	SkipFiles []string
	SkipDirs  []string
}

func (m SyncManager) Run() error {
	if m.SSHClientConfig == nil {
		m.SSHClientConfig = &ssh.ClientConfig{}
	}

	comps, err := url.Parse(m.Remote)
	if err != nil {
		log.ERR.Fatalln("could not parse remote url:", err)
	}

	if comps.User != nil {
		m.SSHClientConfig.User = comps.User.Username()

		if pwd, ok := comps.User.Password(); ok {
			m.SSHClientConfig.Auth = append(m.SSHClientConfig.Auth, ssh.Password(pwd))
		} else {
			host, port, err := net.SplitHostPort(comps.Host)
			if err != nil {
				log.WRN.Println("could not extract host,port:", err)
			} else {
				nPort, err := strconv.Atoi(port)
				if err != nil {
					log.WRN.Println("could not resolve port:", err)
				} else {
					pwd, err := seckeychain.FindInternetPassword(host, "", m.SSHClientConfig.User, "", uint16(nPort), seckeychain.ProtocolTypeSSH, seckeychain.AuthenticationTypeAny)
					if err != nil {
						println("ssh://" + comps.Host)
						log.WRN.Println("could not access keychain:", err)
					} else {
						m.SSHClientConfig.Auth = append(m.SSHClientConfig.Auth, ssh.Password(pwd))
					}
				}
			}
		}
	}

	log.NFO.Println("Dialing", comps)
	client, err := internal.NewSFTPClient(comps.Host, m.SSHClientConfig)
	if err != nil {
		log.ERR.Fatalln("could not connect sftp:", err)
	}

	log.NFO.Println("Listing files")
	remoteMap := internal.BuildRemoteFileList(client.Walk(comps.Path), comps.Path)

	done := make(chan bool)

	rmdirs, rms, mkdirs, puts := internal.CompareTree(m.Local, remoteMap, m.SkipFiles, m.SkipDirs)
	rmdirc := internal.FileSliceToChan(rmdirs)
	rmc := internal.FileSliceToChan(rms)
	mkdirc := internal.FileSliceToChan(mkdirs)
	putc := internal.FileSliceToChan(puts)

	clients := internal.BuildClients(comps.Host, m.SSHClientConfig, 4)

	for i := 0; i < 4; i++ {
		go internal.ProcessDelete(clients[i], rmc, m.DryRun, done)
	}

	for i := 0; i < 4; i++ {
		<-done
	}

	for i := 0; i < 1; i++ {
		go internal.ProcessDelete(clients[i], rmdirc, m.DryRun, done)
	}

	for i := 0; i < 1; i++ {
		<-done
	}

	for i := 0; i < 1; i++ {
		go internal.ProcessPut(clients[i], comps.Path, mkdirc, m.DryRun, done)
	}

	for i := 0; i < 1; i++ {
		<-done
	}

	for i := 0; i < 4; i++ {
		go internal.ProcessPut(clients[i], comps.Path, putc, m.DryRun, done)
	}

	for i := 0; i < 4; i++ {
		<-done
	}

	return nil
}
