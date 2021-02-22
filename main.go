package main

import (
    "log"
    "net/http"
    
    c "./controllers"
)

const (
    PORT = ":8080"
)

func main () {
    log.Println("FILE MANAGER")
    http.HandleFunc("/files",c.RetrieveFiles)
    http.HandleFunc("/upload",c.UploadFile)
    http.HandleFunc("/delete",c.DeleteFiles)
    http.HandleFunc("/update",c.Rename)
    http.HandleFunc("/dir",c.NewFolder)
    http.HandleFunc("/assets", c.ServeImage)
    log.Fatal(http.ListenAndServe(PORT, nil))
}