package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Client struct{}

func NewClient() (*Client, error) {
	return &Client{}, nil
}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (app *Client) ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	const maxBytes = 1024 * 1024 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

func (app *Client) WriteJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (app *Client) ErrorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var jsonResponse JsonResponse
	jsonResponse.Error = true
	jsonResponse.Message = err.Error()

	return app.WriteJson(w, statusCode, jsonResponse)
}
