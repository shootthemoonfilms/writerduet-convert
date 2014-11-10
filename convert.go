package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	format = flag.String("format", "fountain", "Target format (fountain, celtx, fdx)")
	script = flag.String("script", "", "Source script")
	dest   = flag.String("dest", "", "Destination file")
)

func main() {
	flag.Parse()

	switch *format {
	case "fountain":
		break
	case "celtx":
		break
	case "fdx":
		break
	default:
		flag.PrintDefaults()
		fmt.Println("Valid formats are fountain, celtx, or fdx")
		return
	}

	if *script == "" || *dest == "" {
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(*script)
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("script", filepath.Base(*script))
	if err != nil {
		panic(err)
		return
	}
	_, err = io.Copy(part, file)

	// Write other supporting fields
	_ = writer.WriteField("MAX_FILE_SIZE", "100000000")
	_ = writer.WriteField("format", *format)
	_ = writer.WriteField("agree_tou", "on")

	err = writer.Close()
	if err != nil {
		panic(err)
		return
	}

	req, err := http.NewRequest("POST", "https://writerduet.com/convert_script.php", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Referer", "https://writerduet.com/convert")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		return
	}

	body = &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
		return
	}
	resp.Body.Close()
	err = ioutil.WriteFile(*dest, body.Bytes(), 0644)
	if err != nil {
		panic(err)
		return
	}
}

// POST https://writerduet.com/convert_script.php
// Referer: https://writerduet.com/convert
// Content-Type: multipart/form-data

// Payload:

// Content-Disposition: form-data; name="MAX_FILE_SIZE"
//
// 100000000

// Content-Disposition: form-data; name="script"; filename="SCRIPT.celtx"
// Content-Type: application/octet-stream

// Content-Disposition: form-data; name="format"
//
// fountain

// (other options are "fdx" (Final Draft) and "celtx" (Celtx))

// Content-Disposition: form-data; name="agree_tou"
//
// on
