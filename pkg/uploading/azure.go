package uploading

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func AzureUpload(filepath *string, account *string, container *string, gzipUpload bool) {
	accountKey := os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountKey) == 0 {
		log.Fatal("AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}
	credential, err := azblob.NewSharedKeyCredential(*account, accountKey)
	check(err)
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", *account, *container))
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	file, err := os.Open(*filepath)
	check(err)
	reader, writer := io.Pipe()
	go compress(writer, file, gzipUpload)
	fileNameBase := path.Base(file.Name())
	if gzipUpload {
		fileNameBase = fileNameBase + ".gz"
	}
	ctx := context.Background()
	fmt.Printf("Uploading the file with blob name: %s\n", *filepath)
	_, err = azblob.UploadStreamToBlockBlob(ctx, reader, containerURL.NewBlockBlobURL(fileNameBase), azblob.UploadStreamToBlockBlobOptions{
		BufferSize: 5 * 1024 * 1024,
		MaxBuffers: 10})
	check(err)
	fmt.Printf("done")
}
