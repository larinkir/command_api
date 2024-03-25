package server

import (
	"fmt"
	"net/http"
	"os"
)

func Start() error {
	const op = "http-server.server.Start"
	port := os.Getenv("PORT")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		return fmt.Errorf("%s. Failed start server:\n%w", op, err)
	}
	return nil
}
