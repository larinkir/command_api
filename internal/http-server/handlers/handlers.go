package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/larinkir/command_api/internal/bash"
	"github.com/larinkir/command_api/internal/cash"
	"github.com/larinkir/command_api/internal/http-server/response"
	"github.com/larinkir/command_api/internal/storage/psql"
	"log"
	"net/http"
)

type Request struct {
	CommandName string `json:"command_name,omitempty"`
	Parameters  string `json:"parameters,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Storage
type Storage interface {
	SaveCommand(commandName string) (*psql.ReturnCommand, error)
	GetAllCommands() ([]psql.Command, error)
	GetCommand(id string) (*psql.Command, error)
	DeleteCommand(id int, name string) (*psql.ReturnCommand, error)
}

func HandleSaver(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.Saver"
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("ERROR: %s. Empty body request.", op)
				response.Error(http.StatusBadRequest, "empty body request", w)
				return
			} else {
				log.Printf("ERROR: %s. Error decoding request:\n%v", op, err)
				response.Error(http.StatusBadRequest, "request decoding failed", w)
				return
			}
		}

		if req.CommandName == "" {
			log.Printf("ERROR: %s. Empty request received.", op)
			response.Error(http.StatusBadRequest, "empty request", w)
			return
		}
		com, err := storage.SaveCommand(req.CommandName)
		if err != nil {
			if err.Error() == "duplicate command" {
				response.Error(http.StatusBadRequest, "duplicate command", w)
				return
			}
			response.Error(http.StatusInternalServerError, err.Error(), w)
			return
		}

		log.Printf("OK: %s. Command %s with id %d successfully added.", op, com.Name, com.Id)
		response.OK(com, w)

	}
}

func HandleGetterAllCommands(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.GetterAllCommands"
		commandList, err := storage.GetAllCommands()
		if err != nil {
			response.Error(http.StatusInternalServerError, err.Error(), w)
			return
		}
		log.Printf("OK: %s. All commands successfully getted.", op)
		response.OK(commandList, w)
	}
}

func HandleGetterCommand(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.GetterCommand"
		id := r.URL.Path[8:]
		if id == "" {
			response.Error(http.StatusBadRequest, "empty id", w)
			log.Printf("ERROR: %s. Id is empty", op)
			return
		}

		com, err := storage.GetCommand(id)
		if err != nil {
			response.Error(http.StatusInternalServerError, err.Error(), w)
			return
		}
		log.Printf("OK: %s. Command: %s successfully getted.", op, com.Name)

		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("ERROR: %s. Error decoding the request:\n%v", op, err)
			response.Error(http.StatusBadRequest, "request decoding failed", w)
			return
		}
		execCommand, err := bash.RunCommand(com, req.Parameters)
		if err != nil {
			response.Error(http.StatusInternalServerError, "failed to run the command", w)
			return
		}
		response.OK(execCommand, w)
	}
}

func HandleDeleter(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.DeleterCommand"
		id := r.URL.Path[11:]

		if id == "" {
			response.Error(http.StatusBadRequest, "empty id", w)
			log.Printf("ERROR: %s. Id is empty", op)
			return
		}

		com, err := storage.GetCommand(id)

		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				response.Error(http.StatusBadRequest, "command does not exist", w)
				return
			}
			response.Error(http.StatusInternalServerError, err.Error(), w)
			return
		}
		delCom, err := storage.DeleteCommand(com.Id, com.Name)
		if err != nil {
			response.Error(http.StatusInternalServerError, err.Error(), w)
			return
		}
		log.Printf("OK: %s. Command %s with id %d successfully deleted.", op, com.Name, com.Id)

		response.OK(delCom, w)
	}
}

func HandleStopperCommand(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.StopperCommand"
		id := r.URL.Path[9:]
		com, err := storage.GetCommand(id)
		if err != nil {
			response.Error(http.StatusInternalServerError, "failed to get the command", w)
			return
		}
		if cash.CommandProcess[com.Id] == nil {
			log.Printf("ERROR: command has not been run before")
			response.Error(http.StatusBadRequest, "command has not been run before", w)
			return
		}
		execCommand, err := bash.StopCommand(cash.CommandProcess[com.Id], com)
		if err != nil {
			response.Error(http.StatusInternalServerError, "failed to stop the command", w)
			return
		}
		response.OK(execCommand, w)
	}
}

func FirstPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API is working.")
	}
}
