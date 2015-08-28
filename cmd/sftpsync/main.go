package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iwat/go-log"
	"github.com/iwat/sftpsync"
)

var skipFiles = []string{
	".buildpath",
	".DS_Store",
	".git",
	".gitignore",
	".gitmodules",
	".project",
}

var skipDirs = []string{
	".git",
	".settings",
}

var flagDryRun = false
var flagVerbose = false
var flagVeryVerbose = false

func init() {
	flag.BoolVar(&flagDryRun, "dryrun", false, "Dry-run")
	flag.BoolVar(&flagVerbose, "v", false, "Verbose")
	flag.BoolVar(&flagVeryVerbose, "vv", false, "Very verbose")
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 && flag.NArg() != 2 {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "Usage: %s <remote> [<local>]\n", os.Args[0])
		os.Exit(2)
	}

	log.Init(flagVerbose, flagVeryVerbose)

	local := "."
	if flag.NArg() == 2 {
		local = flag.Arg(1)
	}

	m := sftpsync.SyncManager{
		Local:     local,
		Remote:    flag.Arg(0),
		SkipFiles: skipFiles,
		SkipDirs:  skipDirs,
		DryRun:    flagDryRun,
	}
	err := m.Run()
	if err != nil {
		log.ERR.Fatal(err)
	}
}
