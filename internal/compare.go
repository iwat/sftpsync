package internal

import (
	"path/filepath"
	"regexp"
	"sort"

	"github.com/kr/fs"

	"github.com/iwat/go-log"
)

func CompareTree(basepath string, remoteMap map[string]file, skipFiles, skipDirs []*regexp.Regexp, appendFile bool) (rmdirs []file, rms []file, mkdirs []file, puts []file) {
	walker := fs.Walk(basepath)

	for walker.Step() {
		if err := walker.Err(); err != nil {
			log.WRN.Println("walker error:", err)
			continue
		}

		rel, err := filepath.Rel(basepath, walker.Path())
		if err != nil {
			log.WRN.Println("rel error:", err)
			continue
		}

		if rel == "." {
			continue
		}

		name := filepath.Base(walker.Path())
		if walker.Stat().IsDir() {
			matched := false
			for _, skipDir := range skipDirs {
				if skipDir.MatchString(name) {
					walker.SkipDir()
					matched = true
					break
				}
			}
			if matched {
				continue
			}
		} else {
			matched := false
			for _, skipFile := range skipFiles {
				if skipFile.MatchString(name) {
					matched = true
					break
				}
			}
			if matched {
				continue
			}
		}

		stat := walker.Stat()
		mine := file{
			mode:    stat.Mode(),
			size:    stat.Size(),
			mod:     stat.ModTime(),
			path:    walker.Path(),
			relPath: rel,
		}

		if remote, ok := remoteMap[mine.relPath]; ok {
			if !mine.mode.IsDir() && !remote.mode.IsDir() {
				if mine.mod.After(remote.mod) {
					puts = append(puts, mine)
				} else if mine.size != remote.size {
					if remote.size < mine.size && appendFile {
						mine.offset = remote.size
					}
					puts = append(puts, mine)
				}
			} else if !mine.mode.IsDir() && remote.mode.IsDir() {
				rmdirs = append(rmdirs, remote)
				puts = append(puts, mine)
			} else if mine.mode.IsDir() && !remote.mode.IsDir() {
				rms = append(rms, remote)
			}
			delete(remoteMap, remote.relPath)
		} else {
			if mine.mode.IsDir() {
				mkdirs = append(mkdirs, mine)
			} else {
				puts = append(puts, mine)
			}
		}
	}

	for _, remote := range remoteMap {
		if remote.mode.IsDir() {
			rmdirs = append(rmdirs, remote)
		} else {
			rms = append(rms, remote)
		}
	}

	sort.Sort(fileByPath(mkdirs))
	sort.Sort(fileByPath(puts))
	sort.Reverse(fileByPath(rmdirs))
	sort.Reverse(fileByPath(rms))

	return rmdirs, rms, mkdirs, puts
}
