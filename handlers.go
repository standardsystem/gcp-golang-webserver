package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
)

const largeFileThreshold = 30 * 1024 * 1024

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	userEmail := GetUserEmail(ctx)
	bucketname := os.Getenv("BUCKET")

	if len(bucketname) == 0 {
		log.Println("| 500 |", userEmail, r.Host, r.Method, r.URL.Path)
		http.Error(w, "環境変数 BUCKETが設定されていません。", 500)
		return
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println("| 500 |", userEmail, err.Error(), r.Host, r.Method, r.URL.Path)
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

	// if file not found
	// try to find index.html
	if err != nil {
		filename = path.Join(filename, "index.html")

		oh = bucket.Object(filename)
		objAttrs, err = oh.Attrs(ctx)
	}

	if err != nil {
		elapsed := time.Since(start)
		log.Println("| 404 |", userEmail, elapsed.String(), r.Host, r.Method, r.URL.Path, err.Error())
		http.Error(w, "Not found", 404)
		return
	}

	if objAttrs.Size > largeFileThreshold {
		url, err := bucket.SignedURL(filename, &storage.SignedURLOptions{
			Scheme:  storage.SigningSchemeV4,
			Method:  "GET",
			Expires: time.Now().Add(15 * time.Minute),
		})
		if err == nil {
			log.Println("| 302 |", userEmail, "Redirecting large file to Signed URL", filename, objAttrs.Size)
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		log.Println("| 500 |", userEmail, "Failed to generate Signed URL", err.Error())
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
		log.Println("| 200 |", userEmail, elapsed.String(), err.Error(), r.Host, r.Method, r.URL.Path)
		return
	}

	elapsed := time.Since(start)
	log.Println("| 200 |", userEmail, elapsed.String(), r.Host, r.Method, r.URL.Path)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	log.Println("|", err)
	return err == nil
}
