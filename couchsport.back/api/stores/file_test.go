package stores

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/amaurybrisou/couchsport.back/api/types"
	"github.com/golang/leveldb/db"
	"github.com/golang/leveldb/memfs"
)

type memFS struct {
	_os db.FileSystem
}

func (mem memFS) OpenFile(name string) (io.WriteCloser, error) {
	f, err := mem._os.Create(name)
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

func TestFileStore_Save(t *testing.T) {
	type fields struct {
		FileSystem    types.FileSystem
		PublicPath    string
		ImageBasePath string
		FilePrefix    string
	}
	type args struct {
		directory string
		filename  string
		buf       io.Reader
	}

	memos := memFS{_os: memfs.New()}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should return an error",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    "public/",
				ImageBasePath: "static/img/",
				FilePrefix:    "isupload.",
			},
			args: args{
				directory: "user-",
				filename:  "test-file.jpg",
				buf:       strings.NewReader(``),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "should return correct filename",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    "public/",
				ImageBasePath: "static/img",
				FilePrefix:    "isupload.",
			},
			args: args{
				directory: "user-3",
				filename:  "test-file-1.jpg",
				buf:       strings.NewReader(`tototototo`),
			},
			want:    "static/img/user-3/isupload.test-file-1.jpg",
			wantErr: false,
		},
		{
			name: "empty prefix",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    "public/dist/",
				ImageBasePath: "static/img/",
				FilePrefix:    "isupload.",
			},
			args: args{
				directory: "3",
				filename:  "test-file-1.jpg",
				buf:       strings.NewReader(`tototototo`),
			},
			want:    "static/img/3/isupload.test-file-1.jpg",
			wantErr: false,
		},
		{
			name: "empty prefix",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    "public/dist/",
				ImageBasePath: "static/img/",
				FilePrefix:    "isupload.",
			},
			args: args{
				directory: "3",
				filename:  "test-file-1.jpg",
				buf:       strings.NewReader(`tototototo`),
			},
			want:    "static/img/3/isupload.test-file-1.jpg",
			wantErr: false,
		},
		{
			name: "empty file prefix",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    "public/",
				ImageBasePath: "static/img/",
				FilePrefix:    "",
			},
			args: args{
				directory: "user-3",
				filename:  "test-file-1.jpg",
				buf:       strings.NewReader(`tototototo`),
			},
			want:    "static/img/user-3/test-file-1.jpg",
			wantErr: false,
		},
		{
			name: "empty filename",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    "public/",
				ImageBasePath: "static/img/",
				FilePrefix:    "isupload.",
			},
			args: args{
				directory: "user-",
				filename:  "",
				buf:       strings.NewReader(`tototototo`),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty filename",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    "public/",
				ImageBasePath: "static/img/",
				FilePrefix:    "",
			},
			args: args{
				directory: "",
				filename:  "toto.jpg",
				buf:       strings.NewReader(`tototototo`),
			},
			want:    "static/img/toto.jpg",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fileStore{
				FileSystem:    tt.fields.FileSystem,
				PublicPath:    tt.fields.PublicPath,
				ImageBasePath: tt.fields.ImageBasePath,
				FilePrefix:    tt.fields.FilePrefix,
			}
			got, err := app.Save(tt.args.directory, tt.args.filename, tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("FileStore.Save() = have %v, want %v", got, tt.want)
				return
			}

			if tt.wantErr {
				return
			}

			tmp, err := memos.Stat(tt.fields.PublicPath + got)
			if err != nil {
				t.Errorf("couldn stat file %s", got)
				return
			}
			t.Logf("file correctly stat %s", tmp)
		})
	}
}
