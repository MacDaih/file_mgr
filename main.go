package main

import (
    "log"
    "net/http"

    c "./controllers"
)

func main () {
    log.Println("FILE MANAGER")

    http.HandleFunc("/files",c.RetrieveFiles)
    http.HandleFunc("/upload",c.UploadFile)
    http.HandleFunc("/delete",c.DeleteFiles)
    http.HandleFunc("/update",c.Rename)
    http.HandleFunc("/dir",c.NewFolder)

    log.Fatal(http.ListenAndServe(":8000", nil))
}