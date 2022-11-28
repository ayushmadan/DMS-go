package main

import (
	// "errors"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	azureblob()
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	azureblob()

}

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}

func azureblob() {
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

func GetMessage() {
	run := true
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "20.219.165.69:9093",
		"group.id":          "foo",
		"auto.offset.reset": "smallest"})
	topics := []string{}
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.SubscribeTopics(topics, nil)

	for run == true {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Fprintf(os.Stdin, "%% Error: %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}

}
func htmltopdf() {
	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	pdfg.Grayscale.Set(true)

	// Create a new input page from an URL
	page := wkhtmltopdf.NewPage("./index.html")

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile("./simplesample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
