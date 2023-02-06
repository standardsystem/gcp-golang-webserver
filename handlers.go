package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := context.Background()
	bucketname := os.Getenv("BUCKET")

	if len(bucketname) == 0 {
		log.Println("| 500 |", r.Host, r.Method, r.URL.Path)
		http.Error(w, "環境変数 BUCKETが設定されていません。", 500)
		return
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println("| 500 |", err.Error(), r.Host, r.Method, r.URL.Path)
		http.Error(w, "ストレージクライアントが作成できません。", 500)
		return
	}
	bucket := client.Bucket(bucketname)
	filename := r.URL.Path[1:]

	if filename == "/" || filename == "" {
		filename = "index.html"
	}

	oh := bucket.Object(filename)
	objAttrs, err := oh.Attrs(ctx)
	if err != nil {
		elapsed := time.Since(start)
		log.Println("| 404 |", elapsed.String(), r.Host, r.Method, r.URL.Path, err.Error())
		http.Error(w, "Not found", 404)
		return
	}

	o := oh.ReadCompressed(true)
	rc, err := o.NewReader(ctx)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", objAttrs.ContentType)
	w.Header().Set("Content-Encoding", objAttrs.ContentEncoding)
	w.Header().Set("Content-Length", strconv.Itoa(int(objAttrs.Size)))
	w.WriteHeader(200)

	if _, err := io.Copy(w, rc); err != nil {
		elapsed := time.Since(start)
		log.Println("| 200 |", elapsed.String(), err.Error(), r.Host, r.Method, r.URL.Path)
		return
	}

	elapsed := time.Since(start)
	log.Println("| 200 |", elapsed.String(), r.Host, r.Method, r.URL.Path)
}


func fileExists(filename string) bool {
    _, err := os.Stat(filename)
	log.Println("|", err)
    return err == nil
}