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
	"github.com/jung-kurt/gofpdf"
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
		pdfPath := tempFile.Name() + ".pdf"
		err = markdownToPDF(tempFile.Name(), pdfPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cmd = exec.Command("mv", pdfPath, outputPath)
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

	tmpl, err := template.New("").Parse(`<h2>Conversion Result:</h2><pre>{{.}}</pre><a href="/download/{{.}}">Download Result</a>`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, result)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["filename"]
	http.ServeFile(w, r, fileName)
}

func markdownToPDF(inputPath, outputPath string) error {
		pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	_, err = pdf.MultiCell(0, 10, string(content), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	if err != nil {
		return err
	}

	return pdf.OutputFileAndClose(outputPath)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/download/{filename}", downloadHandler) // New route for downloading

	http.Handle("/", r)

	port := "8080"
	fmt.Printf("Server listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
