package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	medium "github.com/medium/medium-sdk-go"
)

var (
	flagFile   = flag.String("file", "", "Required: Filename with Markdown or image")
	flagToken  = flag.String("secret", "", "Required: Medium Integration Token")
	flagType   = flag.String("type", "blog", "Type of upload: blog or image")
	flagTitle  = flag.String("title", "no title given", "Title of blog post")
	flagStatus = flag.String("status", "draft", "Status of post: draft, published or unlisted")
	flagTags   = flag.String("tags", "", "Comma seperated line of blog tags")
)

// GetFileContentType from https://golangcode.com/get-the-content-type-of-file/
func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func UploadBlog() {
	fmt.Println("Uploading blog....")

	md, err := ioutil.ReadFile(*flagFile)
	if err != nil {
		log.Fatal(err)
	}

	m := medium.NewClientWithAccessToken(*flagToken)

	u, err := m.GetUser("")
	if err != nil {
		log.Fatal(err)
	}

	status := medium.PublishStatusDraft
	if *flagStatus == "published" {
		status = medium.PublishStatusPublic
	}

	p, err := m.CreatePost(medium.CreatePostOptions{
		UserID:        u.ID,
		Title:         *flagTitle,
		Tags:          strings.Split(*flagTags, ","),
		Content:       string(md),
		ContentFormat: medium.ContentFormatMarkdown,
		PublishStatus: status,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Blog post URL: %s", p.URL)
}

func uploadImage() {
	fmt.Println("Uploading image....")

	// Open File
	f, err := os.Open(*flagFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Get the content
	contentType, err := GetFileContentType(f)
	if err != nil {
		panic(err)
	}

	fmt.Println("Content Type: " + contentType)

	m := medium.NewClientWithAccessToken(*flagToken)

	u, err := m.GetUser("")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Uploading to user: %v", u)

	i, err := m.UploadImage(medium.UploadOptions{
		FilePath:    *flagFile,
		ContentType: "",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Image URL: %s", i.URL)
}
func main() {
	flag.Parse()

	if *flagFile == "" {
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Please go to https://medium.com/me/settings and generate an Integration Tokens.")
		fmt.Println("If that setting category is not available email yourfriends@medium.com and request access.")
		fmt.Println()
		log.Fatal("Error: --file and --secret are required")
	}

	if *flagType == "blog" {
		UploadBlog()
	}

	if *flagType == "image" {
		uploadImage()
	}

}
