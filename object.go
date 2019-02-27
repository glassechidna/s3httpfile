package s3httpfile

import (
	"bytes"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"time"
)

type s3file struct {
	input  *s3.GetObjectInput
	output *s3.GetObjectOutput
	reader *bytes.Reader
}

func (f *s3file) Name() string {
	return *f.input.Key
}

func (f *s3file) Size() int64 {
	return *f.output.ContentLength
}

func (f *s3file) Mode() os.FileMode {
	return 0400
}

func (f *s3file) ModTime() time.Time {
	return *f.output.LastModified
}

func (f *s3file) IsDir() bool {
	return false
}

func (f *s3file) Sys() interface{} {
	return nil
}

func (f *s3file) Close() error {
	f.reader = nil
	return nil
}

func (f *s3file) Read(p []byte) (n int, err error) {
	return f.reader.Read(p)
}

func (f *s3file) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

func (f *s3file) Readdir(count int) ([]os.FileInfo, error) {
	panic("not a dir")
}

func (f *s3file) Stat() (os.FileInfo, error) {
	return f, nil
}
