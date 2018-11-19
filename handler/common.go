package handler

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chinx/cobweb"
)

func urlParam(r *http.Request, key string) string {
	params := r.Context().Value(context.Background())
	if params == nil {
		return ""
	}
	return params.(cobweb.Params).Get(key)
}

func readBody(body io.Reader, v interface{}) error {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, v)
	if err != nil {
		return err
	}
	return nil
}

func reply(w http.ResponseWriter, status int, data interface{}, err error) {
	var result []byte
	switch t := data.(type) {
	case []byte:
		result = t
	case string:
		result = []byte(t)
	case error:
		result = []byte(t.Error())
	default:
		byteData, err := json.Marshal(data)
		if err != nil {
			status = http.StatusInternalServerError
			result = []byte(err.Error())
		} else {
			result = byteData
		}
	}
	if status >= http.StatusBadRequest {
		if err != nil {
			log.Println(err)
		} else {
			log.Println(status, string(result))
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(result)
}

