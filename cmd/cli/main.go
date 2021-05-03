package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/kacurez/keboola-sdk-go/pkg/uploading"
)

type UploadCommand struct {
	CloudProvider              string `arg:"positional,required"`
	File                       string `arg:"positional,required"`
	DestinationBucketOrAccount string `arg:"positional,required"`
	KeyOrContainerPath         string `arg:"-k" help:"s3 key path" default:""`
	Gzip                       bool   `arg:"-g, --gzip" help:"gzip on upload"`
}

func doUpload(args *UploadCommand) {
	switch args.CloudProvider {
	case "S3":
		uploading.S3Upload(&args.File, &args.DestinationBucketOrAccount, &args.KeyOrContainerPath, args.Gzip)
	case "Azure":
		uploading.AzureUpload(&args.File, &args.DestinationBucketOrAccount, &args.KeyOrContainerPath, args.Gzip)
	default:
		fmt.Println("unknown provider:" + args.CloudProvider)
	}
}

func main() {
	var args struct {
		Upload *UploadCommand `arg:"subcommand:upload"`
	}
	p := arg.MustParse(&args)
	switch {
	case args.Upload != nil:
		fmt.Println("upload")
		doUpload(args.Upload)

	default:
		p.WriteHelp(os.Stdout)
	}

}
