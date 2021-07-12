package util

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

type Resource struct {
	Link     string
	FileName string
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 5*time.Second)
}

func download(url, path string) error {
	fmt.Println("Downloading: " + url)

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialTimeout,
		},
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	_, err = io.Copy(out, io.TeeReader(resp.Body, buf))
	if err != nil {
		return err
	}

	return nil
}

func DownloadResources(resources []Resource, path string) error {
	for _, r := range resources {
		attempts := 3

		for attempts > 0 {
			attempts--

			err := download(r.Link, path+"/"+r.FileName)

			if err == nil {
				break
			}
		}

	}

	return nil
}
