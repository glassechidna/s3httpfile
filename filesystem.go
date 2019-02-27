package s3httpfile

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type S3FileSystem struct {
	ShowDirectoryListing bool
	api                  s3iface.S3API
	bucket               string
	prefix               string
}

func NewS3FileSystem(api s3iface.S3API, bucket, prefix string) *S3FileSystem {
	return &S3FileSystem{
		api:    api,
		bucket: bucket,
		prefix: prefix,
	}
}

func (fs *S3FileSystem) openDir(name string) (http.File, error) {
	// if we got to this point, the file was not found. if anyone has a key
	// component named `index.html`, then i guess we will have to revisit this.
	if filepath.Base(name) == "index.html" || !fs.ShowDirectoryListing {
		return nil, os.ErrNotExist
	}

	objects := []*s3.Object{}
	prefixes := []*s3.CommonPrefix{}

	inputPrefix := fs.prefix + name
	if name != "" {
		inputPrefix += "/"
	}

	input := &s3.ListObjectsV2Input{
		Prefix:    &inputPrefix,
		Delimiter: aws.String("/"),
		Bucket:    &fs.bucket,
	}
	err := fs.api.ListObjectsV2Pages(input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		objects = append(objects, page.Contents...)
		prefixes = append(prefixes, page.CommonPrefixes...)
		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return &s3dir{
		input:    input,
		contents: objects,
		prefixes: prefixes,
	}, nil
}

func (fs *S3FileSystem) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		return fs.openDir("")
	}

	input := &s3.GetObjectInput{
		Bucket: &fs.bucket,
		Key:    aws.String(fs.prefix + name),
	}
	output, err := fs.api.GetObject(input)

	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == s3.ErrCodeNoSuchKey {
			return fs.openDir(name)
		}
	}

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(output.Body)
	buf := bytes.NewReader(body)

	return &s3file{
		input:  input,
		output: output,
		reader: buf,
	}, nil

}
