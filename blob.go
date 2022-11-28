package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func uploadToAzureBlob() {
	fmt.Printf("Azure Blob storage quick start sample\n")
	ctx := context.Background()
	connectionString := os.Getenv("AZURE_STORAGE_ACCOUNT_CONNECTION_STRING")
	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	// Create the container
	data := []byte("\nhello world this is a blob\n")
	_, err1 := serviceClient.UploadBuffer(ctx, "test", "test-go", data, nil)
	if err1 != nil {
		log.Fatal(err1)
	}
}
