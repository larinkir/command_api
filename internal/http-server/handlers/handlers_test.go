package handlers

import (
	"bytes"
	"fmt"
	"github.com/larinkir/command_api/internal/http-server/handlers/mocks"
	"github.com/larinkir/command_api/internal/storage/psql"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleSaver(t *testing.T) {
	cases := []struct {
		name          string //Название кейса
		command       string //Имя команды
		errMock       error  //Возможная ошибка из мок функции
		expStatusCode int    //Ожидаемый статус
		expResp       string //Ожидаемый вывод

	}{
		{
			name:          "save_ok",
			command:       "test",
			expStatusCode: 200,
			expResp:       "{\"command\":{\"id\":100,\"name\":\"test\"},\"message\":\"command successfully added\"}",
		},
		{
			name:          "save_empty_request",
			command:       "",
			expStatusCode: 400,
			expResp:       "{\"error\":\"empty request\"}",
		},
		{
			name:          "save_duplicate_command",
			command:       "test",
			errMock:       fmt.Errorf("duplicate command"),
			expStatusCode: 400,
			expResp:       "{\"error\":\"duplicate command\"}",
		},
		{
			name:          "save_unknown_error",
			command:       "test",
			errMock:       fmt.Errorf("unknown error"),
			expStatusCode: 500,
			expResp:       "{\"error\":\"unknown error\"}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)

			if tc.errMock != nil {
				storage.On("SaveCommand", tc.command).
					Return(nil, tc.errMock)
			}

			if !strings.Contains(tc.expResp, "error") {
				storage.On("SaveCommand", tc.command).
					Return(&psql.ReturnCommand{
						Command: psql.Command{
							Id:   100,
							Name: tc.command},
						Message: "command successfully added",
					}, nil)

			}

			handler := HandleSaver(storage)

			input := fmt.Sprintf(`{"command_name": "%s"}`, tc.command)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.expStatusCode)

			resp := rr.Body.String()
			require.Equal(t, tc.expResp, resp)

		})
	}
}

func TestHandleGetterAllCommands(t *testing.T) {
	cases := []struct {
		name          string //Название кейса
		errMock       error  //Возможная ошибка из мок функции
		expStatusCode int    //Ожидаемый статус
		expResp       string //Ожидаемый вывод
	}{
		{
			name:          "getall_ok",
			expStatusCode: 200,
			expResp:       "{\"command\":{},\"command_list\":[{\"id\":101,\"name\":\"test1\"},{\"id\":102,\"name\":\"test2\"},{\"id\":103,\"name\":\"test3\"}]}",
		},
		{
			name:          "getall_error_db",
			errMock:       fmt.Errorf("unknown error"),
			expStatusCode: 500,
			expResp:       "{\"error\":\"unknown error\"}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)

			if tc.errMock != nil {
				storage.On("GetAllCommands").
					Return(nil, tc.errMock)
			} else {
				storage.On("GetAllCommands").
					Return([]psql.Command{{Id: 101, Name: "test1"}, {Id: 102, Name: "test2"}, {Id: 103, Name: "test3"}}, nil)
			}

			handler := HandleGetterAllCommands(storage)

			req, err := http.NewRequest(http.MethodGet, "/getall", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.expStatusCode)

			resp := rr.Body.String()
			require.Equal(t, tc.expResp, resp)

		})
	}
}

func TestHandleDeleter(t *testing.T) {
	cases := []struct {
		name          string //Название кейса
		id            int    // Id удаляемой команды
		command       string //Название команды
		errMock       error  //Возможная ошибка из мок функции
		expStatusCode int    //Ожидаемый статус
		expResp       string //Ожидаемый вывод
	}{
		{
			name:          "delete_command_not_exist",
			id:            100,
			command:       "test",
			expStatusCode: 400,
			expResp:       "{\"error\":\"command does not exist\"}",
		},
		{
			name:          "delete_ok",
			id:            100,
			command:       "test",
			expStatusCode: 200,
			expResp:       "{\"command\":{\"id\":100,\"name\":\"test\"},\"message\":\"command successfully deleted\"}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)

			if tc.errMock != nil {
				storage.On("DeleteCommand", tc.id, tc.command).
					Return(nil, tc.errMock)
			}

			if !strings.Contains(tc.expResp, "error") {
				storage.On("GetCommand", fmt.Sprintf("%d", tc.id)).
					Return(&psql.Command{
						Id:   tc.id,
						Name: tc.command,
					}, nil)
				storage.On("DeleteCommand", tc.id, tc.command).
					Return(&psql.ReturnCommand{
						Command: psql.Command{
							Id:   tc.id,
							Name: tc.command},
						Message: "command successfully deleted",
					}, nil)

				handler := HandleDeleter(storage)

				req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/delete/id=%d", tc.id), nil)
				require.NoError(t, err)

				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				require.Equal(t, rr.Code, tc.expStatusCode)

				resp := rr.Body.String()
				require.Equal(t, tc.expResp, resp)
			}
		})
	}
}
