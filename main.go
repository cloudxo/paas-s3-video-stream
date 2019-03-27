package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	bucketName, ok := os.LookupEnv("S3_BUCKET_NAME")
	if !ok {
		log.Fatalf("S3_BUCKET_NAME must be set")
	}
	bucketKey, ok := os.LookupEnv("S3_BUCKET_KEY")
	if !ok {
		log.Fatalf("S3_BUCKET_KEY must be set")
	}
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	filename := "s3-download"
	f, err := ioutil.TempFile("", "paas-s3-video-stream-download-*")
	if err != nil {
		log.Fatalf("failed to create file %q, %v", filename, err)
	}

	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
	})
	if err != nil {
		log.Fatalf("failed to download file, %v", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, f.Name())
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
