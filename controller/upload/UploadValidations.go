package upload

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path"
	"slices"

	"github.com/isaquecsilva/static-server/utils"
)

func generateUploadValidations() utils.ValidationActions {
	return map[string]func(args ...any) (error, int) {
		"filename": func(args ...any) (error, int) {
			header := args[0].(*multipart.FileHeader)

			if header.Filename == "" {
				return errors.New("invalid filename"), http.StatusBadRequest
			}

			return nil, -1
		},

		"extension": func(args ...any) (error, int) {
			header := args[0].(*multipart.FileHeader)
			fileTypeList := args[1].([]string)

			ext := path.Ext(header.Filename)

			if !slices.Contains(fileTypeList, ext) {
				return errors.New("file type not allowed"), http.StatusBadRequest
			}

			return nil, -1
		},

		"size": func(args ...any) (error, int) {
			header  := args[0].(*multipart.FileHeader)

			if args[1].(int64) < header.Size {
				return fmt.Errorf("file size overtake the maximum allowed: %v", args[2]), http.StatusBadRequest
			}

			return nil, -1
		},
	}
}