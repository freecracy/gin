package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

func Test_dedupEnv(t *testing.T) {
	type args struct {
		env []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test_dedupEnv",
			args: args{
				[]string{"a=a", "b=b", "a=c"},
			},
			want: []string{"a=c", "b=b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dedupEnv(tt.args.env); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dedupEnv() = %v, want %v", got, tt.want)
				// ???????
				// dedupEnv() = [   a=c b=b], want [a=c b=b]
			}
		})
	}
}
