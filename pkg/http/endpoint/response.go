package endpoint

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"olympos.io/encoding/edn"
)

var (
	ErrResponseReaded    = errors.New("endpoint: response readed")
	ErrUnmappedMediaType = errors.New("endpoint: unmapped media type")
)

type Response interface {
	Unmarshal(dest interface{}) error
	Status() int
	RawBody() (string, []byte, error)
	Headers() http.Header
	Close() error
}

type response struct {
	Response *http.Response
	readed   bool
}

func (r *response) Unmarshal(dest interface{}) error {
	contentType, body, err := r.RawBody()
	if err != nil {
		return err
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	switch mediaType {
	case "application/json":
		return json.Unmarshal(body, dest)
	case "application/edn":
		return edn.Unmarshal(body, dest)
	default:
		return ErrUnmappedMediaType
	}
}

func (r *response) Status() int {
	return r.Response.StatusCode
}

func (r *response) RawBody() (string, []byte, error) {
	if r.readed {
		return "", nil, ErrResponseReaded
	}
	r.readed = true
	contentType := r.Response.Header.Get("content-type")
	body, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return "", nil, err
	}
	return contentType, body, nil
}

func (r *response) Headers() http.Header {
	return r.Response.Header
}

func (r *response) Close() error {
	r.readed = true
	return r.Response.Body.Close()
}
