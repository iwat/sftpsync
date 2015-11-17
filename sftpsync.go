package sftpsync

import (
	"net"
	"net/url"
	"regexp"
	"strconv"

	"golang.org/x/crypto/ssh"

	"github.com/iwat/go-seckeychain"
	"github.com/iwat/sftpsync/internal"
	"gopkg.in/inconshreveable/log15.v2"
)

type SyncManager struct {
	Append bool
	DryRun bool
	Local  string
	Remote string

	SSHClientConfig *ssh.ClientConfig

	SkipFiles []*regexp.Regexp
	SkipDirs  []*regexp.Regexp
}

var log log15.Logger

func init() {
	log = log15.New("package", "sftpsync")
}

func (m SyncManager) Run() error {
	if m.SSHClientConfig == nil {
		m.SSHClientConfig = &ssh.ClientConfig{}
	}

	comps, err := url.Parse(m.Remote)
	if err != nil {
		log.Crit("could not parse remote url", "err", err)
		return err
	}

	if comps.User != nil {
		m.SSHClientConfig.User = comps.User.Username()

		if pwd, ok := comps.User.Password(); ok {
			m.SSHClientConfig.Auth = append(m.SSHClientConfig.Auth, ssh.Password(pwd))
		} else {
			host, port, err := net.SplitHostPort(comps.Host)
			if err != nil {
				log.Warn("could not extract host,port", "err", err)
			} else {
				nPort, err := strconv.Atoi(port)
				if err != nil {
					log.Warn("could not resolve port", "err", err)
				} else {
					pwd, err := seckeychain.FindInternetPassword(host, "", m.SSHClientConfig.User, "", uint16(nPort), seckeychain.ProtocolTypeSSH, seckeychain.AuthenticationTypeAny)
					if err != nil {
						log.Warn("could not access keychain", "err", err)
					} else {
						m.SSHClientConfig.Auth = append(m.SSHClientConfig.Auth, ssh.Password(pwd))
					}
				}
			}
		}
	}

	log.Info("Dialing", "uri", comps)
	client, err := internal.NewSFTPClient(comps.Host, m.SSHClientConfig)
	if err != nil {
		log.Crit("could not connect sftp", "err", err)
		return err
	}

	log.Info("Listing files")
	remoteMap := internal.BuildRemoteFileList(client.Walk(comps.Path), comps.Path)

	done := make(chan bool)

	rmdirs, rms, mkdirs, puts := internal.CompareTree(m.Local, remoteMap, m.SkipFiles, m.SkipDirs, m.Append)
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
