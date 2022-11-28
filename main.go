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
	"github.com/segmentio/kafka-go"
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
	testConn()
	GetMessage()
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

	w := &kafka.Writer{
		Addr:                   kafka.TCP(os.Getenv("KAFKA_CONNECTION_STRING")),
		Topic:                  "test",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	go func() {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{os.Getenv("KAFKA_CONNECTION_STRING")},
			Topic:   "test",
			GroupID: "my-group",
		})
		for {
			// the `ReadMessage` method blocks until we receive the next event
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				panic("could not read message " + err.Error())
			}
			// after receiving the message, log its value
			fmt.Println("received: ", string(msg.Value))
		}
	}()
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages

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
func testConn() {
	conn, err := kafka.Dial("tcp", os.Getenv("KAFKA_CONNECTION_STRING"))
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		fmt.Println(k)
	}
}
