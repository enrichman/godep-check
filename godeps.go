package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Godeps struct {
	ImportPath   string
	GoVersion    string
	GodepVersion string
	Packages     []string `json:",omitempty"`
	Deps         []Dependency
}

type Dependency struct {
	ImportPath string
	Comment    string `json:",omitempty"`
	Rev        string
}

func loadGodeps() (*Godeps, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(dir + "/Godeps.json")
	if err != nil {
		bytes, err = ioutil.ReadFile(dir + "/Godeps/Godeps.json")
		if err != nil {
			return nil, err
		}
	}

	var gdeps Godeps
	err = json.Unmarshal(bytes, &gdeps)
	if err != nil {
		return nil, err
	}

	gdeps.Deps = mergeDeps(gdeps.Deps)

	return &gdeps, nil
}

func mergeDeps(deps []Dependency) []Dependency {
	alreadyLoadedDeps := make(map[string]bool)
	mergedDeps := make([]Dependency, 0)

	for _, d := range deps {
		splitted := strings.Split(d.ImportPath, "/")
		cutLen := 3
		if len(splitted) < 3 {
			cutLen = len(splitted)
		}
		d.ImportPath = strings.Join(splitted[:cutLen], "/")
		found := alreadyLoadedDeps[d.ImportPath]
		if !found {
			alreadyLoadedDeps[d.ImportPath] = true
			mergedDeps = append(mergedDeps, d)
		}
	}

	return mergedDeps
}
