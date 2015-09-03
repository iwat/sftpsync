package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/iwat/go-log"
	"github.com/iwat/sftpsync"
)

var skipFiles = []*regexp.Regexp{}
var skipDirs = []*regexp.Regexp{}

var flagDryRun = false
var flagVerbose = false
var flagVeryVerbose = false

func init() {
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.buildpath$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.DS_Store$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.git$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.gitignore$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.gitmodules$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.project$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^myapp_.*$`))

	skipDirs = append(skipDirs, regexp.MustCompile(`^\.git$`))
	skipDirs = append(skipDirs, regexp.MustCompile(`^\.settings$`))

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
