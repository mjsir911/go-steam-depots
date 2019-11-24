package main

import (
	"./steam"
	"archive/zip"
	"fmt"
	// "bytes"
	// "net/http"
	"os"
)


func main() {
	// r, _ := http.Get("http://cache20-iad1.steamcontent.com/depot/881101/manifest/2114008103684897778/5")
	// buf := make([]byte, r.ContentLength)
	// io.ReadFull(r.Body, buf)
	// d := bytes.NewReader(buf)
	// unzipper, _ := zip.NewReader(d, int64(d.Len()))

	// d, _ := os.Open("5")
	// info, _ := d.Stat()
	// unzipper, _ := zip.NewReader(d, info.Size())
	// f, _ := unzipper.File[0].Open()
	// var a steam.Depot
	// a, _ = steam.NewDepot(f)
	// var b steam.File
	// b = a.Files[0]
	// // a = steam.Depot{}
	// fmt.Println(b)

	f, _ := os.Open("test.json")
}
