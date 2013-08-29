package main

import (
	"archive/zip"
	"log"
	"path"
	"runtime"
	"sort"
	"sync"
)

var timeLocs []string
var tlOnce sync.Once

func listTimeLocations() ([]string, error) {
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
