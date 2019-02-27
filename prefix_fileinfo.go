package s3httpfile

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"path/filepath"
	"time"
)

type s3PrefixFileInfo struct {
	*s3.CommonPrefix
}

func (fi *s3PrefixFileInfo) Name() string {
	return filepath.Base(*fi.CommonPrefix.Prefix)
}

func (fi *s3PrefixFileInfo) Size() int64 {
	return -1
}

func (fi *s3PrefixFileInfo) Mode() os.FileMode {
	return 0400
}

func (fi *s3PrefixFileInfo) ModTime() time.Time {
	return time.Now()
}

func (fi *s3PrefixFileInfo) IsDir() bool {
	return true
}

func (fi *s3PrefixFileInfo) Sys() interface{} {
	return nil
}
