package main

import (
	"archive/zip"
	"log"
	"os"
	"path"
	"runtime"
	"sort"
	"sync"
	"time"
)

var timeLocs []string
var tlOnce sync.Once

func findfiles(p, prefix string, files []string) []string {
	d, err := os.Open(p)
	if err != nil {
		return files
	}
	defer d.Close()

	infos, err := d.Readdir(-1)
	if err != nil {
		return files
	}

	for _, info := range infos {
		if info.Mode().IsRegular() {
			files = append(files, prefix+info.Name())
		} else if info.IsDir() {
			files = findfiles(path.Join(p, info.Name()), info.Name()+"/", files)
		}
	}

	return files
}

func listTimeLocations() ([]string, error) {
	for _, p := range []string{"/usr/share/zoneinfo", "/usr/share/lib/zoneinfo", "/usr/lib/locale/TZ"} {
		files := findfiles(p, "", nil)
		duprem := make(map[string]bool)
		for _, loc := range files {
			if _, err := time.LoadLocation(loc); err == nil {
				duprem[loc] = true
			}
		}
		var locs []string
		for loc := range duprem {
			locs = append(locs, loc)
		}
		if len(locs) > 0 {
			sort.Strings(locs)
			return locs, nil
		}
	}

	zoneinfoZip := path.Join(runtime.GOROOT(), "lib", "time", "zoneinfo.zip")
	z, err := zip.OpenReader(zoneinfoZip)
	if err != nil {
		return nil, err
	}
	defer z.Close()

	locs := []string{}
	for _, f := range z.File {
		if f.Name[len(f.Name)-1] == '/' {
			continue
		}
		locs = append(locs, f.Name)
	}

	sort.Strings(locs)
	return locs, nil
}

func loadTimeLocs() {
	tlOnce.Do(func() {
		var err error
		if timeLocs, err = listTimeLocations(); err != nil {
			log.Fatalf("Could not load time locations: %s", err)
		}
	})
}
