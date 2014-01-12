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
	FORMAT = flag.String("format", "fountain", "Target format (fountain, celtx, fdx)")
	SCRIPT = flag.String("script", "", "Source script")
	DEST   = flag.String("dest", "", "Destination file")
)

func main() {
	flag.Parse()

	switch *FORMAT {
	case "fountain":
		break
	case "celtx":
		break
	case "fdx":
		break
	default:
		fmt.Printf("Valid formats are fountain, celtx, or fdx\n")
		return
	}

	if *SCRIPT == "" {
		flag.PrintDefaults()
		return
	}
	if *DEST == "" {
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(*SCRIPT)
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("script", filepath.Base(*SCRIPT))
	if err != nil {
		panic(err)
		return
	}
	_, err = io.Copy(part, file)

	// Write other supporting fields
	_ = writer.WriteField("MAX_FILE_SIZE", "100000000")
	_ = writer.WriteField("format", *FORMAT)
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
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			panic(err)
			return
		}
		resp.Body.Close()
		err = ioutil.WriteFile(*DEST, body.Bytes(), 0644)
		if err != nil {
			panic(err)
			return
		}
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
