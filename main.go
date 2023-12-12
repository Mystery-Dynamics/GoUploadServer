package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// const uploadDir = "static/uploads" // Directory to store uploaded files
var uploadDir string

func main() {
	secretKey := flag.String("secret", "your_secret_key", "Secret key for upload")
	pathUpload := flag.String("pathUpload", "static/uploads", "Path to upload")
	pageUpload := flag.String("pageUpload", "/upload-page", "Path to upload")
	port := flag.Int("port", 8088, "Port to listen on")
	enableUploadPage := flag.Bool("enableUploadPage", false, "Port to listen on")
	allowDownload := flag.Bool("allowDownload", true, "Port to listen on")
	flag.Parse()
	uploadDir = *pathUpload
	// Create upload directory if it doesn't exist
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		fmt.Printf("Error creating upload directory: %v", err)
		return
	}
	if *enableUploadPage {
		http.HandleFunc(*pageUpload, uploadPageHandler)
	}
	if *allowDownload {
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(uploadDir)))) // Serve static files
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler(*secretKey))

	fmt.Printf("Server listening on port %d\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func uploadPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/upload.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func uploadHandler(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Verify secret key
		if r.FormValue("secret") != secret {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Invalid secret key")
			return
		}

		// Get uploaded file
		file, header, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error retrieving uploaded file: %v", err)
			return
		}
		defer file.Close()

		// Read file data
		data, err := ioutil.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error reading uploaded file: %v", err)
			return
		}

		// Generate unique filename
		filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), header.Filename)

		// Save uploaded file
		filepath := filepath.Join(uploadDir, filename)
		err = ioutil.WriteFile(filepath, data, 0644)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error saving uploaded file: %v", err)
			return
		}

		downloadURL := fmt.Sprintf("/static/%s", filename)

		// Generate download link HTML
		downloadLink := fmt.Sprintf("<p>Download file: </p><a href='%s'>%s</a>", downloadURL, downloadURL)

		// Set content type to HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Respond with success message and download link
		fmt.Fprintf(w, "<!DOCTYPE html><html><body><h1>File uploaded successfully: %s</h1><p>%s</p></body></html>", filename, downloadLink)
	}
}
