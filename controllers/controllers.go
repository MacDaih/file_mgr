package controllers

import (
	"log"
    "os"
	"io"
	"fmt"
    "io/ioutil"
    "net/http"
    "path/filepath"
    "encoding/json"

	u "../utils"
	m "../models"
    l "../logs"
)

const PATH = "./tmp-rep/"

func RetrieveFiles(w http.ResponseWriter, r *http.Request) {
    u.SetCors(&w, "GET")
	var path string
	keys, ok := r.URL.Query()["sub"]
    if !ok || len(keys[0]) < 1 {
        path = PATH
    } else {
		path = fmt.Sprintf("%s%s",PATH,keys[0])
	}

    files, err := ioutil.ReadDir(path)
    
    if err != nil {
        l.WriteLog("RetrieveFiles",nil, err)
        http.Error(w,fmt.Sprintf("%v",err),http.StatusInternalServerError)
        return
    }

    var res []m.File

    for _,f := range files {
        id,name := u.GetFileID(f.Name())
        if len(*name) == 0 {
            name = id
        }
        file := m.File{
            *id,
            *name,
            f.Size(),
            f.ModTime(),
            f.IsDir(),
        }        
        res = append(res,file)
    }

    repo := m.Repo{
        path,
        res,
    }
    json.NewEncoder(w).Encode(repo)
} 

func UploadFile(w http.ResponseWriter, r *http.Request) {
    u.SetCors(&w, "POST")
    r.ParseMultipartForm(10 << 20)
	var id string
    file, handler, err := r.FormFile("upload")
	
    if err != nil {
        ft := fmt.Sprintf("%v",err)
        http.Error(w, ft, http.StatusBadRequest)
        return
    }
    defer file.Close()
	if id, err = u.NewFileId(); err != nil {
		http.Error(w,fmt.Sprintf("%v",err),http.StatusInternalServerError)
        return
	}
    noun := u.FormatFileDST(id,handler.Filename)
	fullPath := fmt.Sprintf("%s%s",PATH,noun)
	dst, err := os.Create(fullPath)
	defer dst.Close()

    if err != nil {
        l.WriteLog("UploadFile",nil, err)
        http.Error(w,"Failed",http.StatusInternalServerError)
        return
	}

	if _, err := io.Copy(dst, file); err != nil {
        l.WriteLog("UploadFile",nil, err)
		http.Error(w, "Failed", http.StatusInternalServerError)
		return
	}
    l.WriteLog("UploadFile",&fullPath, nil)
    w.WriteHeader(http.StatusOK)
}

func DeleteFiles(w http.ResponseWriter, r *http.Request) {
    u.SetCors(&w, "POST")
    var d []m.File

    b, err := ioutil.ReadAll(r.Body)

    if err != nil {
        l.WriteLog("UploadFile",nil, err)
        http.Error(w,"Failed",http.StatusInternalServerError)
        return
    }

    err = json.Unmarshal(b, &d)
    
    if err != nil {
        l.WriteLog("DeleteFiles",nil, err)
        http.Error(w,"Failed",http.StatusInternalServerError)
        return
    }
    for _, f := range d {
        rebuild := fmt.Sprintf("%s%s_%s",PATH,f.Id,f.Name)
        err := os.RemoveAll(rebuild)
        if err != nil {
            l.WriteLog("DeleteFiles",nil, err)
            http.Error(w,"Failed",http.StatusInternalServerError)
            return
        }
    }
}

func NewFolder(w http.ResponseWriter, r *http.Request) {
    u.SetCors(&w, "POST")

    body, err := ioutil.ReadAll(r.Body)

    if err != nil {
        log.Println(err)
        http.Error(w,"Failed",http.StatusInternalServerError)
        return 
    }

    data := struct{
        DirName string `json:"dirName"`
    }{}

    err = json.Unmarshal(body, &data)
    if err != nil {
        l.WriteLog("NewFolder",nil, err)
        http.Error(w,"Failed",http.StatusInternalServerError)
        return
    }
	id, err := u.NewFileId()
	if err != nil {
        l.WriteLog("NewFolder",nil, err)
		http.Error(w,"Failed",http.StatusInternalServerError)
        return
	}
	fullName := fmt.Sprintf("%s_%s",id, data.DirName)
    full := fmt.Sprintf("%s%s", PATH, fullName)
    err = os.Mkdir(full, 0755)
    if err != nil {
        l.WriteLog("NewFolder",nil, err)
        http.Error(w,"Failed",http.StatusInternalServerError)
        return
    }
    l.WriteLog("NewFolder",&full, nil)
	w.WriteHeader(http.StatusOK)
}

func Rename(w http.ResponseWriter, r *http.Request) {
	u.SetCors(&w, "POST")
	body, err := ioutil.ReadAll(r.Body)

    if err != nil {
        l.WriteLog("Rename",nil, err)
        http.Error(w,"Failed",http.StatusInternalServerError)
        return 
    }
    data := struct{
		PrevName string `json:"prev_name"`
        NewName string `json:"new_name"`
    }{}

    err = json.Unmarshal(body, &data)
    if err != nil {
        l.WriteLog("Rename",nil, err)
        http.Error(w,"Failed",http.StatusInternalServerError)
        return
    }
	src := fmt.Sprintf("%s%s",PATH, data.PrevName)
	dst := fmt.Sprintf("%s%s",PATH, data.NewName)
    err = os.Rename(src, dst)
	if err != nil {
        l.WriteLog("Rename",nil, err)
        http.Error(w,"Failed",http.StatusInternalServerError)
        return
	}
}
// WIP !!!
func ServeImage(w http.ResponseWriter, r *http.Request) {
    u.SetCors(&w, "GET")
    _, ok := r.URL.Query()["id"]
    if !ok {
        http.Error(w,"Failed",http.StatusBadRequest)
    }
    err := filepath.Walk(PATH,
        func(path string, info os.FileInfo, err error) error {
            if err != nil {
                l.WriteLog("ServeImage",nil, err)
                return err
            }
            fmt.Println(path, info.Size())
            u.DeepDown(path)
            return nil
        })
    if err != nil {
        l.WriteLog("ServeImage",nil, err)
    }
    w.WriteHeader(http.StatusOK)
}