package bash

import (
	"bufio"
	"github.com/larinkir/command_api/internal/cash"
	"github.com/larinkir/command_api/internal/storage/psql"
	"io"
	"log"
	"os/exec"
	"strings"
)

type ExecCommand struct {
	Id      int
	Name    string
	Output  []string
	Message string
}

func RunCommand(com *psql.Command, param string) (*ExecCommand, error) {
	const op = "internal.bash.RunCommand"
	var cmd *exec.Cmd
	if param == "" {
		cmd = exec.Command(com.Name)
	} else {
		cmd = exec.Command(com.Name, param)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("ERROR: %s. Error receiving the command output stream: %w.", op, err)
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		log.Printf("ERROR: %s. Error starting the command: %w.", op, err)
		return nil, err
	}

	cash.CommandProcess[com.Id] = cmd

	reader := bufio.NewReader(stdout)
	var outputList []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("ERROR: %s. Error reading a string from the command output stream: %w.", op, err)
			return nil, err
		}
		line = strings.TrimSuffix(line, "\n")
		outputList = append(outputList, line)
	}

	err = cmd.Wait()
	if err != nil {
		if err.Error() == "signal: killed" {
			delete(cash.CommandProcess, com.Id)
			return &ExecCommand{Id: com.Id,
				Name:    com.Name,
				Output:  outputList,
				Message: "The command successfully stopped",
			}, nil
		} else {
			log.Printf("ERROR: %s. Error waiting for command completion: %w.", op, err)
		}
		return nil, err
	}
	delete(cash.CommandProcess, com.Id)

	log.Printf("OK: %s. Command %s successfully executed.", op, com.Name)

	return &ExecCommand{Id: com.Id,
		Name:   com.Name,
		Output: outputList,
	}, nil
}

func StopCommand(cmd *exec.Cmd, com *psql.Command) (*ExecCommand, error) {
	const op = "internal.bash.StopCommand"

	err := cmd.Process.Kill()
	if err != nil {
		log.Printf("ERROR: %s. Error stopping command: %w", op, err)
		return nil, err
	}
	log.Printf("OK: %s. The command %s was successfully stopped.", op, com.Name)
	return &ExecCommand{Id: com.Id,
		Name:    com.Name,
		Message: "The command successfully stopped",
	}, nil

}
