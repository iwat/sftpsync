package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/iwat/sftpsync"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/inconshreveable/log15.v2/ext"
)

var skipFiles = []*regexp.Regexp{}
var skipDirs = []*regexp.Regexp{}

var log log15.Logger

var (
	flagAppend      = false
	flagDryRun      = false
	flagVerbose     = false
	flagVeryVerbose = false
)

func init() {
	log = log15.New("package", "main")

	skipFiles = append(skipFiles, regexp.MustCompile(`^\.buildpath$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.DS_Store$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.git$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.gitignore$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.gitmodules$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^\.project$`))
	skipFiles = append(skipFiles, regexp.MustCompile(`^myapp_.*$`))

	skipDirs = append(skipDirs, regexp.MustCompile(`^\.git$`))
	skipDirs = append(skipDirs, regexp.MustCompile(`^\.settings$`))

	flag.BoolVar(&flagAppend, "append", false, "Append instead of overwrite")
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

	handler := ext.FatalHandler(log15.CallerFileHandler(log15.StderrHandler))
	if flagVeryVerbose {
		handler = log15.LvlFilterHandler(log15.LvlDebug, handler)
	} else if flagVerbose {
		handler = log15.LvlFilterHandler(log15.LvlInfo, handler)
	}
	log15.Root().SetHandler(handler)

	local := "."
	if flag.NArg() == 2 {
		local = flag.Arg(1)
	}

	m := sftpsync.SyncManager{
		Local:     local,
		Remote:    flag.Arg(0),
		SkipFiles: skipFiles,
		SkipDirs:  skipDirs,
		Append:    flagAppend,
		DryRun:    flagDryRun,
	}
	err := m.Run()
	if err != nil {
		log.Crit("run error", "err", err)
		return
	}
}
