package main

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_unpackZip(t *testing.T) {
	type args struct {
		targetDir   string
		archiveFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "dl",
			args: args{
				targetDir:   filepath.Join(os.Getenv("HOME"), ".gvm2", "dl"),
				archiveFile: filepath.Join(os.Getenv("HOME"), ".gvm2", "dl", "dl.zip"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := unpackZip(tt.args.targetDir, tt.args.archiveFile); (err != nil) != tt.wantErr {
				t.Errorf("unpackZip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
