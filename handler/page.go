package handler

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/chinx/wepub/api"
	"github.com/chinx/wepub/module"
)

func CreatePage(w http.ResponseWriter, r *http.Request) {
	result := checkAdmin(w, r)
	if result.Status != module.StatusLogin {
		return
	}
	pageURL := fmt.Sprintf("pages/%d.html", time.Now().UnixNano())
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		reply(w, http.StatusInternalServerError, err, nil)
		return
	}
	file, err := os.Create(StaticDir + "/" + pageURL)
	if err != nil {
		reply(w, http.StatusInternalServerError, err, nil)
		return
	}

	api.HTMLTemplate.Execute(file, map[string]interface{}{
		"title":    "发条互动",
		"htmlBody": template.HTML(content),
	})
	file.Close()

	fullURL := "http://" + r.Host + "/editor/" + pageURL
	if err != nil {
		reply(w, http.StatusInternalServerError, err, nil)
		return
	}
	reply(w, http.StatusCreated, fullURL, nil)
}

func GetAction(w http.ResponseWriter, r *http.Request) {
	result := checkAdmin(w, r)
	if result.Status != module.StatusLogin {
		return
	}
	switch r.URL.Query()["action"][0] {
	case "config":
		reply(w, http.StatusOK, api.ConfigJSON, nil)
	default:
		reply(w, http.StatusBadRequest, "action is error", nil)
	}
}

func PostAction(w http.ResponseWriter, r *http.Request) {
	result := checkAdmin(w, r)
	if result.Status != module.StatusLogin {
		return
	}
	switch r.URL.Query()["action"][0] {
	case "uploadimage", "uploadfile", "uploadvideo":
		uploadFile(w, r)
	case "catchimage":
		catchImage(w, r)
	default:
		reply(w, http.StatusBadRequest, "action is error", nil)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	number := r.FormValue("number")
	name := r.FormValue("name")
	err := os.MkdirAll(".\\attachment\\"+number+name, 0777)
	if err != nil {
		reply(w, http.StatusInternalServerError, err, nil)
		return
	}

	file, header, err := r.FormFile("upfile")
	if err != nil {
		reply(w, http.StatusInternalServerError, err, nil)
		return
	}
	defer file.Close()

	path1 := ".\\attachment\\" + number + name + "\\" + header.Filename

	outFile, err := os.Create(path.Join(StaticDir, "upload", path1))
	if err != nil {
		reply(w, http.StatusInternalServerError, err, nil)
		return
	}

	defer outFile.Close()
	io.Copy(outFile, file)

	reply(w, http.StatusOK, map[string]interface{}{
		"state":    "SUCCESS",
		"url":      "/attachment/" + number + name + "/" + header.Filename,
		"title":    header.Filename,
		"original": header.Filename,
	}, nil)
}

func catchImage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	files, ok := r.Form["source[]"]
	if !ok {
		reply(w, http.StatusBadRequest, "error from source", nil)
		return
	}
	list := make([]*api.ImageData, 0, len(files))
	for i := range files {
		img := api.NewImageData()
		img.Source = files[i]
		list = append(list, img)
	}

	result := map[string]interface{}{
		"state": "SUCCESS",
	}
	var err error
	list, err = downSave(list)
	if err != nil {
		result["state"] = "FAILURE"
	}
	result["list"] = list

	reply(w, http.StatusOK, result, nil)
}

type urlChan struct {
	index int
	filename string
}

func downSave(list []*api.ImageData) ([]*api.ImageData, error) {
	var nErr error

	c := make(chan*urlChan)
	for i := range list {
		_, err := url.Parse(list[i].Source)
		if err != nil {
			list[i].State = "FAILURE"
			nErr = err
			continue
		}

		go downSaveOne(list[i].Source, i, c)

		resp, err := http.Get(list[i].Source)
		if err != nil || resp.StatusCode >= http.StatusBadRequest {
			list[i].State = "FAILURE"
			nErr = err
			continue
		}
		defer resp.Body.Close()

		filename, err := module.SaveScale(resp.Body, StaticDir+"/upload/", 640)
		if err != nil {
			list[i].State = "FAILURE"
			nErr = err
			continue
		}
		list[i].URL = strings.Replace(filename, StaticDir, "/editor", 1)
	}
	for j := 0 ; j < len(list); j ++ {
		select {
		case data := <- c:
			if data.filename == ""{
				list[data.index].State = "FAILURE"
			}else{
				list[data.index].URL = data.filename
			}
		case <- time.After(time.Second * 30):
			return list, errors.New("timeout")
		}
	}
	log.Println(nErr)
	return list, nErr
}

func downSaveOne(source string, index int, c chan*urlChan)  {
	resp, err := http.Get(source)
	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		log.Println(err)
		c <- &urlChan{index:index}
		return
	}
	defer resp.Body.Close()

	filename, err := module.SaveScale(resp.Body, StaticDir+"/upload/", 640)
	if err != nil {
		log.Println(err)
		c <- &urlChan{index:index}
		return
	}

	c <- &urlChan{index:index, filename: strings.Replace(filename, StaticDir, "/editor", 1)}
}
