package response

import (
	"encoding/json"
	"github.com/larinkir/command_api/internal/bash"
	"github.com/larinkir/command_api/internal/storage/psql"
	"log"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseOK struct {
	Command     psql.Command   `json:"command,omitempty"`
	CommandList []psql.Command `json:"command_list,omitempty"`
	Output      []string       `json:"output,omitempty"`
	Message     string         `json:"message,omitempty"`
}

func Error(code int, error string, w http.ResponseWriter) {
	const op = "internal.http-server.response.Error"
	w.Header().Set("Content-Type", "application/json")
	js, _ := json.Marshal(ResponseError{
		Error: error,
	})
	w.WriteHeader(code)
	_, err := w.Write(js)
	if err != nil {
		log.Printf("ERROR: %s. %v", op, err.Error())
	}
}

func OK(v interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	const op = "internal.http-server.response.OK"
	switch value := v.(type) {
	case []psql.Command:
		js, _ := json.Marshal(ResponseOK{
			CommandList: value,
		})
		_, err := w.Write(js)
		if err != nil {
			log.Printf("ERROR: %s. %v", op, err.Error())
		}

	case *bash.ExecCommand:
		js, _ := json.Marshal(ResponseOK{
			Command: psql.Command{
				Id:   value.Id,
				Name: value.Name,
			},
			Output:  value.Output,
			Message: value.Message,
		})
		_, err := w.Write(js)
		if err != nil {
			log.Printf("ERROR: %s. %v", op, err.Error())
		}
	case *psql.ReturnCommand:
		js, _ := json.Marshal(ResponseOK{
			Command: psql.Command{
				Id:   value.Id,
				Name: value.Name,
			},
			Message: value.Message,
		})
		_, err := w.Write(js)
		if err != nil {
			log.Printf("ERROR: %s. %v", op, err.Error())
		}

	}
}
