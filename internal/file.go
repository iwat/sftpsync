// sftpsync - Sync local file system to SFTP

// Copyright (c) 2015 Chaiwat Shuetrakoonpaiboon. All rights reserved.
//
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

package internal

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type file struct {
	mode    os.FileMode
	size    int64
	mod     time.Time
	path    string
	relPath string
	offset  int64
}

func (f file) String() string {
	return fmt.Sprintf("%s %6d %s %s", f.mode, f.size, f.mod.Format("2006-Jan-02 15:04:05"), f.relPath)
}

type fileByPath []file

func (f fileByPath) Len() int {
	return len(f)
}

func (f fileByPath) Less(a, b int) bool {
	arr := []file(f)
	if !arr[a].mode.IsDir() && arr[b].mode.IsDir() {
		return true
	} else if arr[a].mode.IsDir() && !arr[b].mode.IsDir() {
		return false
	}
	return strings.Compare(arr[a].relPath, arr[b].relPath) < 0
}

func (f fileByPath) Swap(a, b int) {
	arr := []file(f)
	arr[a], arr[b] = arr[b], arr[a]
}

func FileSliceToChan(in []file) <-chan file {
	out := make(chan file)
	go func() {
		for _, e := range in {
			out <- e
		}
		close(out)
	}()

	return out
}
