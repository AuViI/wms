package simpleuser

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/AuViI/wms/model"
	yaml "gopkg.in/yaml.v2"
)

type (
	Request struct {
		Type     string `json:"type"`
		SetUsers Users  `json:"users"`
	}
	Response struct {
		Error   bool   `json:"error"`
		Message string `json:"msg"`
	}
	Users []User
	User  struct {
		ID    uint64      `json:"id" yaml:"id"`
		Name  string      `json:"name" yaml:"name"`
		Theme model.Theme `json:"theme" yaml:"theme"`
	}
)

var (
	rw       sync.RWMutex
	userFile = func() string {
		return path.Join(os.Getenv("HOME"), ".wmsuser.yaml")
	}()
)

func HandleJS(w http.ResponseWriter, r *http.Request) {
	var req Request
	json.NewDecoder(r.Body).Decode(&req)
	writer := json.NewEncoder(w)
	switch req.Type {
	case "set":
		if err := save(userFile, req.SetUsers); err != nil {
			writer.Encode(errorResponse(err))
		}
		return
	case "get":
		u, err := load(userFile)
		if err != nil {
			writer.Encode(errorResponse(err))
		} else {
			writer.Encode(u)
		}
		return
	case "default":
		var user User
		user.ID = 0
		user.Name = "Default"
		user.Theme = model.GetDefaultTheme()
		writer.Encode(user)
		return
	default:
		writer.Encode(errorResponse(errors.New("No known type")))
		return
	}
}

func errorResponse(err error) Response {
	return Response{Error: true, Message: err.Error()}
}

func load(file string) (u Users, err error) {
	rw.RLock()
	defer rw.RUnlock()
	f, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0644)
	err = yaml.NewDecoder(f).Decode(&u)
	return
}

func save(file string, users Users) (err error) {
	rw.Lock()
	defer rw.Unlock()
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	enc := yaml.NewEncoder(f)
	defer enc.Close()
	err = enc.Encode(users)
	return
}
