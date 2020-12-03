package utils

import (
	"io"
	"os"
	"testing"

	"github.com/amaurybrisou/couchsport.back/api/types"
	"github.com/golang/leveldb/db"
	"github.com/golang/leveldb/memfs"
)

type memFS struct {
	_os db.FileSystem
}

func (mem memFS) OpenFile(name string) (io.WriteCloser, error) {
	f, err := mem._os.Open(name)
	return f, err
}

func (mem memFS) Stat(name string) (os.FileInfo, error) {
	return mem._os.Stat(name)
}

func (mem memFS) MkdirAll(path string) error {
	return mem._os.MkdirAll(path, 0700)
}

func (mem memFS) IsNotExist(err error) bool {
	return err != nil
}

func TestCreateDirIfNotExists(t *testing.T) {
	type args struct {
		_os  types.FileSystem
		path string
	}

	memos := memFS{_os: memfs.New()}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "should not return an error and create dir",
			args:    args{_os: memos, path: "/test"},
			want:    "/test",
			wantErr: false,
		},
		{
			name:    "should return /test and no error",
			args:    args{_os: memos, path: "/test"},
			want:    "/test",
			wantErr: false,
		},
		{
			name:    "should return error",
			args:    args{_os: memos, path: ""},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDirIfNotExists(tt.args._os, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CeateDirIfNotExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CeateDirIfNotExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
