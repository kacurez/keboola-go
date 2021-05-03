package uploading

import (
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go/aws"

	// "github.com/aws/aws-sdk-go/aws/credentials"
	"compress/gzip"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var testBucket string = "kacurez"

var testFilename = "test/100mb"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func compress(writer *io.PipeWriter, file *os.File, gzipUpload bool) {
	defer writer.Close()
	dstWriter := io.Writer(writer)
	if gzipUpload {
		gw := gzip.NewWriter(writer)
		defer gw.Close()
		dstWriter = gw
	}
	_, err := io.Copy(dstWriter, file)
	check(err)
}

func S3Upload(filepath *string, bucket *string, key *string, gzipUpload bool) {
	sess := session.Must(session.NewSession(&aws.Config{
		//Credentials: credentials.NewStaticCredentials("asdasd", "asdasd", "aaa"),
		Region: aws.String("us-east-1"),
	}))

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.Concurrency = 10
		u.PartSize = 1024 * 1024 * 5
	})
	file, err := os.Open(*filepath)
	check(err)
	reader, writer := io.Pipe()
	go compress(writer, file, gzipUpload)
	fileNameBase := path.Base(file.Name())
	if gzipUpload {
		fileNameBase = fileNameBase + ".gz"
	}
	s3Path := *key + "/" + fileNameBase
	fmt.Println("Start uploading " + fileNameBase + " to " + *bucket)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: bucket,
		Key:    &s3Path,
		Body:   reader,
	})
	check(err)
	file.Close()
	fmt.Println("done")
}
