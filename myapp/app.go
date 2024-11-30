package myapp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type User struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type fooHandler struct{}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.CreatedAt = time.Now()

	data, _ := json.Marshal(user)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(data))
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func fileServerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving static file:", r.URL.Path)
	http.StripPrefix("/public", http.FileServer(http.Dir("public"))).ServeHTTP(w, r)
}

func uploadsFileServerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving uploaded file:", r.URL.Path)
	http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads"))).ServeHTTP(w, r)
}

// upload handles file upload logic: reading the file, creating directories, and saving the file.
func upload(r *http.Request, uploadDir string) (string, error) {
	// Parse the multipart form
	file, header, err := r.FormFile("upload_file")
	if err != nil {
		return "", fmt.Errorf("failed to get file: %w", err)
	}
	defer file.Close()

	// Ensure the upload directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Save the file to the upload directory
	dstPath := filepath.Join(uploadDir, header.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Return the path to the saved file
	return dstPath, nil
}

// uploadHandler is the HTTP handler for file uploads.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Use the upload function to handle the file upload
	uploadDir := "./uploads"
	savedFilePath, err := upload(r, uploadDir)
	if err != nil {
		http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	fmt.Fprintf(w, "File uploaded successfully: %s", savedFilePath)
}

func NewHttpHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/bar", barHandler)

	mux.Handle("/foo", &fooHandler{})

	mux.HandleFunc("/public/", fileServerHandler)

	// File server for /uploads/
	mux.HandleFunc("/uploads/", uploadsFileServerHandler)

	// File upload endpoint
	mux.HandleFunc("/upload", uploadHandler)

	return mux
}
