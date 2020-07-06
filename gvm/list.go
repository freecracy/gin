package main

import (
	"os"
	"strings"
)

func listAllVersion(targetDir string) ([]string, error) {
	ans := make([]string, 0)
	dir, err := os.Open(targetDir)
	if err != nil {
		return nil, err
	}
	f, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, v := range f {
		if strings.Contains(v.Name(), "beta") {
			continue
		}
		if strings.Contains(v.Name(), "rc") {
			continue
		}
		if !strings.HasPrefix(v.Name(), "go") {
			continue
		}
		if strings.Count(v.Name(), ".") == 2 {
			continue
		}
		if v.Name() == "gotip" || v.Name() == "go.mod" {
			continue
		}
		ans = append(ans, v.Name())
	}
	return ans, nil
}
