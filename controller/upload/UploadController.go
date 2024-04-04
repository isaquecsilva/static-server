package upload

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/isaquecsilva/static-server/model/rules"
)

const uploadPage string = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{ .PageName }}</title>
</head>
<body>
	<form action="/upload" method="POST" enctype="multipart/form-data">
		<input type="file" name="file" accept="{{ .AllowedTypes }}" />
		<input type="submit" name="upload" />
	</form>
</body>
</html>`

var MaxUploadSize int64

type UploadController struct {
	rules *rules.UploadRules
	templ *template.Template
}

func (uc *UploadController) Upload(w http.ResponseWriter, r *http.Request) {
	reader, header, err := r.FormFile("file")
	if err != nil {
		println("Failure Getting MultiPartReader: ", err.Error())
		http.Error(w, "error on uploading", http.StatusInternalServerError)
		return
	}

	if header.Size > MaxUploadSize {
		http.Error(w, fmt.Sprintf("file's size overtake the allowed limit, which is %s", uc.rules.MaxFileSize), http.StatusUnprocessableEntity)
		return
	} else if len(header.Filename) == 0 {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	file, err := os.Create(header.Filename)
	if err != nil {
		println("could not create upload_file: ", err.Error())
		http.Error(w, "failure uploading", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	log.Printf("UPLOADING: %s, SIZE: %d\n", header.Filename, header.Size)
	io.Copy(file, reader)
}

func (uc *UploadController) UploadPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if err := uc.templ.Execute(w, struct {
		PageName     string
		AllowedTypes string
	}{"File Upload", strings.Join(uc.rules.FileTypes.FileTypeList, ",")}); err != nil {
		println("TEMPLATE ERROR: ", err.Error())
	}
}

func NewUploadController(rules *rules.UploadRules) (*UploadController, error) {
	templ, err := template.New("Upload Page").Parse(uploadPage)
	if err != nil {
		return nil, err
	}

	return &UploadController{
		rules,
		templ,
	}, nil
}