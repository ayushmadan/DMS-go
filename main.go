package main

import (
	// "errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"context"
	"math/rand"
	"strconv"
	"time"
	// "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	// "os"
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
	azureblob()
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




func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}

func azureblob() {
	fmt.Printf("Azure Blob storage quick start sample\n")
	ctx := context.Background()
	connectionString:= "DefaultEndpointsProtocol=https;AccountName=doctorlistingingestionpr;AccountKey=+xu1up2YWmlmAvJmzfEsAbyaB7FUoARys0ZOL/oKQTBAurHhnWdL+Ri0GefqG9kxxizP5wJbLBuNon9wu9BZcQ==;EndpointSuffix=core.windows.net;"

	serviceClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	

	// Create the container
	data := []byte("\nhello world this is a blob\n")
	containerName := fmt.Sprintf("quickstart-%s", randomString())
	fmt.Printf("Creating a container named %s\n", containerName)
	_, err1 := serviceClient.UploadBuffer(ctx, "test", "test-go", data, nil)
	if err1 != nil {
		log.Fatal(err1)
	}
}

// func GetMessage(count int, client *azservicebus.Client) {
// 	receiver, err := client.NewReceiverForQueue("myqueue", nil) //Change myqueue to env var
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer receiver.Close(context.TODO())

// 	messages, err := receiver.ReceiveMessages(context.TODO(), count, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, message := range messages {
// 		body := message.Body
// 		fmt.Printf("%s\n", string(body))

// 		err = receiver.CompleteMessage(context.TODO(), message, nil)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }