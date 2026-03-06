package main

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/logos-storage/logos-storage-go-bindings/storage"
)

func main() {
	node, err := storage.New(storage.Config{
		BlockRetries: 5,
	})
	if err != nil {
		log.Fatalf("Failed to create Logos Storage node: %v", err)
	}

	version, err := node.Version()
	if err != nil {
		log.Fatalf("Failed to get Logos Storage version: %v", err)
	}
	log.Printf("Logos Storage version: %s", version)

	if err := node.Start(); err != nil {
		log.Fatalf("Failed to start Logos Storage node: %v", err)
	}
	log.Println("Logos Storage node started")

	buf := bytes.NewBuffer([]byte("Hello World!"))
	len := buf.Len()
	cid, err := node.UploadReader(context.Background(), storage.UploadOptions{Filepath: "hello.txt"}, buf)
	if err != nil {
		log.Fatalf("Failed to upload data: %v", err)
	}
	log.Printf("Uploaded data with CID: %s (size: %d bytes)", cid, len)

	f, err := os.Create("hello.txt")
	if err != nil {
		log.Fatal("Failed to create file:", err)
	}
	defer f.Close()

	opt := storage.DownloadStreamOptions{
		Writer: f,
	}

	if err := node.DownloadStream(context.Background(), cid, opt); err != nil {
		log.Fatalf("Failed to download data: %v", err)
	}

	log.Println("Downloaded data to hello.txt")

	// Wait for a SIGINT or SIGTERM signal
	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	// <-ch

	if err := node.Stop(); err != nil {
		log.Fatalf("Failed to stop Logos Storage node: %v", err)
	}
	log.Println("Logos Storage node stopped")

	if err := node.Destroy(); err != nil {
		log.Fatalf("Failed to destroy Logos Storage node: %v", err)
	}
}
