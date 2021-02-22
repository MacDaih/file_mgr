package models

import "time"

const (
    jpg = "jpg"
    png = "png"
)

type File struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Size int64 `json:"size"`
    ModDate time.Time `json:"modDate"`
    IsDir bool `json:"isDir"`
}

type Repo struct {
    Path string `json:"path"`
    Files []File `json:"files"`
}