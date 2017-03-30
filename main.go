package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type images struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type imageList struct {
	Collection []images `json:"data"`
}

func checkDirectoryExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func createDirectory(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func saveFile(fileURL string, path string) error {
	fileName := strings.Split(fileURL, "/")
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(path + fileName[3])
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	galleryPtr := flag.String("gallery", "", "Sets the imgur gallery to pull images from e.g. r/funny")
	outputPathPtr := flag.String("output", "", "Sets the directory to save the images to e.g. C:\\Pictures\\")
	pageNumberPtr := flag.Int("pageNumber", 1, "Sets the gallery page number to pull images from")

	flag.Parse()

	if len([]rune(*galleryPtr)) == 0 {
		log.Fatal("You must provide a gallery to pull images from.")
		return
	}

	if len([]rune(*outputPathPtr)) == 0 {
		log.Fatal("You must provide an output path to save images to.")
		return
	}

	if strings.HasSuffix(*outputPathPtr, "\\") == false {
		*outputPathPtr += "\\"
	}

	fmt.Println("Output:", *outputPathPtr)

	url := fmt.Sprintf("https://api.imgur.com/3/gallery/%s/time/%d", *galleryPtr, *pageNumberPtr)

	pathExists, err := checkDirectoryExists(*outputPathPtr)
	if err != nil {
		fmt.Println("Directory error:", err)
		return
	}

	if pathExists == false {
		err := createDirectory(*outputPathPtr)
		if err != nil {
			log.Fatal("Error creating directory:", err)
			return
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Request:", err)
		return
	}

	client := &http.Client{}

	req.Header.Add("authorization", "Client-ID b68d6fbd258b5ac")
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do:", err)
		return
	}

	defer resp.Body.Close()

	decodeJSON := json.NewDecoder(resp.Body)

	var imagesSlice imageList

	err = decodeJSON.Decode(&imagesSlice)
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < len(imagesSlice.Collection); i++ {
		err := saveFile(imagesSlice.Collection[i].Link, *outputPathPtr)
		if err != nil {
			fmt.Println("Error saving file:", imagesSlice.Collection[i].Link)
			log.Fatal("Error:", err)
			return
		}
	}

	fmt.Println("Files succesfully saved to:", *outputPathPtr)
}
