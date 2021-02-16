package models

import "time"

type File struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Size int64 `json:"size"`
    ModDate time.Time `json:"modDate"`
    IsDir bool `json:"isDir"`
}