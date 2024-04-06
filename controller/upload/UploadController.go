package upload

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/isaquecsilva/static-server/model/rules"
	"github.com/isaquecsilva/static-server/utils"
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
		<input type="file" name="file" accept="{{ .AllowedTypes }}" multiple />
		<input type="submit" name="upload" />
	</form>
</body>
</html>`

var MaxUploadSize int64

type UploadController struct {
	uploadValidations utils.ValidationActions
	rules             *rules.UploadRules
	templ             *template.Template
}

func (uc *UploadController) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1); err != nil {
		http.Error(w, "could not process the request", http.StatusInternalServerError)
		return
	}

	var checkError = func(err error, errorCode int) bool {
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), errorCode)
			return true
		}

		return false
	}

	for _, file := range r.MultipartForm.File["file"] {
		for key, validation := range uc.uploadValidations {
			switch key {
			case "size":
				err, errorCode := validation(file, MaxUploadSize, uc.rules.MaxFileSize)
				if checkError(err, errorCode) {
					return
				}
			case "extension":
				err, errorCode := validation(file, uc.rules.FileTypes.FileTypeList)
				if checkError(err, errorCode) {
					return
				}
			default:
				err, errorCode := validation(file)
				if checkError(err, errorCode) {
					return
				}
			}
		}

		reader, err := file.Open()
		if err != nil {
			println("could not get multipart file stream: ", err.Error())
			http.Error(w, fmt.Sprintf("the server was unable to get upload file's stream. Please, try again later."), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(os.Stdout, "Uploading %s...\n", file.Filename)
		file, err := os.Create(file.Filename)
		if err != nil {
			println("could not create upload_file: ", err.Error())
			http.Error(w, "failure uploading", http.StatusInternalServerError)
			return
		}

		io.Copy(file, reader)
		reader.Close()
		file.Close()
	}

	w.Write([]byte("Enviado."))
}

func (uc *UploadController) UploadPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	var upPageData struct {
		PageName     string
		AllowedTypes string
	}

	upPageData = struct {
		PageName     string
		AllowedTypes string
	}{
		PageName:     "File Upload",
		AllowedTypes: strings.Join(uc.rules.FileTypes.FileTypeList, ","),
	}

	if err := uc.templ.Execute(w, upPageData); err != nil {
		println("TEMPLATE ERROR: ", err.Error())
	}
}

func NewUploadController(rules *rules.UploadRules) (*UploadController, error) {
	templ, err := template.New("Upload Page").Parse(uploadPage)
	if err != nil {
		return nil, err
	}



	return &UploadController{
		generateUploadValidations(),
		rules,
		templ,
	}, nil
}
