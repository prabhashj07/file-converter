package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	format := r.FormValue("format")

	tempFile, err := os.CreateTemp("", "uploaded-*.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outputPath := tempFile.Name() + "." + format
	var cmd *exec.Cmd

	switch format {
	case "markdown":
		cmd = exec.Command("./file-converter", "convert", tempFile.Name(), outputPath, format)
	case "html":
		cmd = exec.Command("./file-converter", "convert", tempFile.Name(), outputPath, format)
	case "pdf":
		cmd = exec.Command("pdftohtml", tempFile.Name(), outputPath)
	default:
		http.Error(w, "Unsupported format", http.StatusBadRequest)
		return
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := strings.TrimSpace(string(output))

	tmpl, err := template.New("").Parse(`<h2>Conversion Result:</h2><pre>{{.}}</pre>`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, result)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/upload", uploadHandler).Methods("POST")

	http.Handle("/", r)

	port := "8080"
	fmt.Printf("Server listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
