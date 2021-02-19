package utils

import (
	"os"
	"fmt"
	"strings"
	"net/http"

	"github.com/gofrs/uuid"
)

func SetCors(w *http.ResponseWriter, method string) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", method)
} 

func FormatFileDST(id string,inp string) string {
	splt := strings.Split(inp, ".")
	res := fmt.Sprintf("%s_%s.%s",id,splt[0],splt[1])
	return res
}

func GetFileID(full string) (*string, *string) {
	splt := strings.Split(full, "_")
	parsed := ""
	for i,n := range splt {
		if i > 0 {
			if i > 1 {
				parsed += fmt.Sprintf("_%s",n)
			} else {
				parsed += n
			}
		}
	}
	return &splt[0],&parsed
} 

func NewFileId() (string, error) {
	newId, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	value := fmt.Sprintf("%v", newId)
	return value, nil
}

func GetFileType(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	buffer := make([]byte,512)
	_, err = file.Read(buffer)

	if err != nil {
		return err
	}
	filetype := http.DetectContentType(buffer)
	fmt.Println("File type or MIME  : ", filetype)
	return nil
}