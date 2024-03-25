package psql

import (
	"database/sql"
	"fmt"
	"github.com/larinkir/command_api/internal/config"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	db *sql.DB
}

type Command struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ReturnCommand struct {
	Command `json:"command"`
	Message string `json:"message"`
}

func Connect(sd config.StorageData) (*Storage, error) {
	const op = "internal.storage.Connect"
	dataSourceName := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		sd.User, sd.Password, sd.DBName)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("%s. Failed connect to DB:\n%s", op, err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s. Error pinging database:\n%s", op, err.Error())
	}

	log.Println("OK: Database connection established.")

	return &Storage{db: db}, nil
}

func (s *Storage) SaveCommand(commandName string) (*ReturnCommand, error) {
	const op = "internal.storage.SaveCommand"
	var count int
	err := s.db.QueryRow("SELECT COUNT(command) FROM command_list WHERE command=$1", commandName).Scan(&count)
	if count == 1 {
		log.Printf("ERROR: %s. Duplicate command.", op)
		return nil, fmt.Errorf("duplicate command")
	}
	var com Command
	res, err := s.db.Query("INSERT INTO command_list(command) VALUES ($1) RETURNING id, command", commandName)
	if err != nil {
		log.Printf("ERROR: %s. Failed to added command to DB:\n%s", op, err.Error())
		return nil, err
	}

	for res.Next() {
		err := res.Scan(&com.Id, &com.Name)
		if err != nil {
			log.Printf("ERROR: %s. Failed to scan result:\n%s", op, err.Error())
			return nil, err
		}
	}
	return &ReturnCommand{
		Command: com,
		Message: "command successfully added",
	}, nil
}

func (s *Storage) GetAllCommands() ([]Command, error) {
	const op = "internal.storage.GetAllCommands"
	var com Command
	res, err := s.db.Query("SELECT id,command FROM command_list")
	if err != nil {
		log.Printf("ERROR: %s. Failed to get commands from DB:\n%s", op, err.Error())
		return nil, err
	}
	var commandList []Command
	for res.Next() {
		err := res.Scan(&com.Id, &com.Name)
		if err != nil {
			log.Printf("ERROR: %s. Failed to scan result from DB:\n%s", op, err.Error())
			return nil, err
		}
		commandList = append(commandList, com)
	}
	return commandList, nil
}

func (s *Storage) GetCommand(id string) (*Command, error) {
	const op = "internal.storage.GetCommand"
	var com Command
	err := s.db.QueryRow("SELECT id,command FROM command_list WHERE id=$1", id).Scan(&com.Id, &com.Name)
	if err != nil {
		log.Printf("ERROR: %s. Failed to get the command by the id from DB:\n%s", op, err.Error())
		return nil, err
	}
	return &com, nil
}

func (s *Storage) DeleteCommand(id int, name string) (*ReturnCommand, error) {
	const op = "internal.storage.DeleteCommand"
	_, err := s.db.Query("DELETE FROM command_list WHERE id=$1", id)
	if err != nil {
		log.Printf("ERROR: %s. Failed to delete the command from DB:\n%s", op, err.Error())
		return nil, err
	}

	return &ReturnCommand{
		Command: Command{
			Id:   id,
			Name: name},
		Message: "command successfully deleted",
	}, nil
}
