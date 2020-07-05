package main

import (
	"fmt"
	"os"
	"strings"
)

func listAllVersion(targetDir string) error {
	dir, err := os.Open(targetDir)
	if err != nil {
		return err
	}
	f, err := dir.Readdir(-1)
	if err != nil {
		return err
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
	}
	return fmt.Errorf("%v", ans)
}
