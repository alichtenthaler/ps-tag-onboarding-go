package response

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendResponseWithData(t *testing.T) {

	type data struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Age       int    `json:"age"`
	}

	testCases := []struct {
		name       string
		data       interface{}
		statusCode int
	}{
		{
			name:       "Send response status 200",
			data:       data{FirstName: "Ann", LastName: "Peterson", Email: "a@p.com", Age: 20},
			statusCode: 200,
		},
		{
			name:       "Send response status 400",
			data:       GenericError{Error: "error 400"},
			statusCode: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			response := httptest.NewRecorder()
			SendResponse(response, tc.statusCode, tc.data)

			dataBytes, err := json.Marshal(tc.data)
			if err != nil {
				tt.Fatal(err)
			}

			assert.Equal(tt, tc.statusCode, response.Code)
			assert.Equal(tt, append(dataBytes, '\n'), response.Body.Bytes())
			assert.Equal(tt, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

func TestSendResponseWithNoData(t *testing.T) {
	response := httptest.NewRecorder()
	status := http.StatusNoContent
	SendResponse(response, status, nil)

	assert.Equal(t, status, response.Code)
	assert.Equal(t, []byte(nil), response.Body.Bytes())
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
}

func TestSendError(t *testing.T) {
	response := httptest.NewRecorder()
	status := http.StatusInternalServerError
	err := errors.New("server error")

	SendError(response, status, err)

	genErr := GenericError{
		Error: err.Error(),
	}

	errBytes, err := json.Marshal(genErr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, status, response.Code)
	assert.Equal(t, append(errBytes, '\n'), response.Body.Bytes())
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
}

func TestSendValidationError(t *testing.T) {

	response := httptest.NewRecorder()
	status := http.StatusBadRequest
	validationErr := ValidationError{
		Error:   "validation error",
		Details: []string{"error 1", "error 2"},
	}

	SendValidationError(response, status, validationErr)

	errBytes, err := json.Marshal(validationErr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, status, response.Code)
	assert.Equal(t, append(errBytes, '\n'), response.Body.Bytes())
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
}