package s3httpfile

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"path/filepath"
	"time"
)

type s3ObjectFileInfo struct {
	*s3.Object
}

func (fi *s3ObjectFileInfo) Name() string {
	return filepath.Base(*fi.Object.Key)
}

func (fi *s3ObjectFileInfo) Size() int64 {
	return *fi.Object.Size
}

func (fi *s3ObjectFileInfo) Mode() os.FileMode {
	return 0400
}

func (fi *s3ObjectFileInfo) ModTime() time.Time {
	return *fi.Object.LastModified
}

func (fi *s3ObjectFileInfo) IsDir() bool {
	return false
}

func (fi *s3ObjectFileInfo) Sys() interface{} {
	return nil
}
