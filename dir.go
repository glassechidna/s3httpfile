package s3httpfile

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"path/filepath"
	"time"
)

type s3dir struct {
	input    *s3.ListObjectsV2Input
	contents []*s3.Object
	prefixes []*s3.CommonPrefix
}

func (d *s3dir) Name() string {
	return filepath.Base(*d.input.Prefix)
}

func (d *s3dir) Size() int64 {
	panic("implement me")
}

func (d *s3dir) Mode() os.FileMode {
	return 0400
}

func (d *s3dir) ModTime() time.Time {
	latest := time.Time{}
	for _, obj := range d.contents {
		if obj.LastModified.After(latest) {
			latest = *obj.LastModified
		}
	}
	return latest
}

func (d *s3dir) IsDir() bool {
	return true
}

func (d *s3dir) Sys() interface{} {
	return nil
}

func (d *s3dir) Close() error {
	return nil
}

func (d *s3dir) Read(p []byte) (n int, err error) {
	panic("is a dir")
}

func (d *s3dir) Seek(offset int64, whence int) (int64, error) {
	panic("is a dir")
}

func (d *s3dir) Readdir(count int) ([]os.FileInfo, error) {
	fis := []os.FileInfo{}

	for _, obj := range d.contents {
		if *obj.Key != *d.input.Prefix {
			fis = append(fis, &s3ObjectFileInfo{obj})
		}
	}

	for _, prefix := range d.prefixes {
		fis = append(fis, &s3PrefixFileInfo{prefix})
	}

	return fis, nil
}

func (d *s3dir) Stat() (os.FileInfo, error) {
	return d, nil
}
