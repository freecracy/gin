package main

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_listAllVersion(t *testing.T) {
	type args struct {
		targetDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				targetDir: filepath.Join(os.Getenv("HOME"), ".gvm2", "dl"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := listAllVersion(tt.args.targetDir); (err != nil) != tt.wantErr {
				t.Errorf("listAllVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
