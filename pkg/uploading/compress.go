package uploading

import (
	"compress/gzip"
	"io"
	"os"
)

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
