package routes

import (
	"net/http"

	"github.com/isaquecsilva/static-server/controller/upload"
	"github.com/isaquecsilva/static-server/middlewares"
)

func InitRoutes(uploadController *upload.UploadController) {
	http.Handle("GET /upload", middlewares.ConnectionLogger(
		http.HandlerFunc(uploadController.UploadPage),
	))
	
	http.Handle("POST /upload", middlewares.ConnectionLogger(
		http.HandlerFunc(uploadController.Upload),
	))
}

func InitDefaultHandler(dir *string) {
	handler := http.FileServer(http.Dir(*dir))
	http.Handle("GET /", handler)
}
