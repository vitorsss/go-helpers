package endpoint

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"olympos.io/encoding/edn"
)

var (
	ErrInvalidUnmarshalTarget = errors.New("endpoint: invalid unmarshal target")
	ErrResponseReaded         = errors.New("endpoint: response readed")
	ErrUnmappedMediaType      = errors.New("endpoint: unmapped media type")
)

type Response interface {
	Unmarshal(dest interface{}) error
	Status() int
	RawBody() (string, []byte, error)
	RawBodyStream() (string, io.Reader, error)
	Headers() http.Header
	Close() error
}

type response struct {
	Response *http.Response
	readed   bool
}

func (r *response) Unmarshal(dest interface{}) error {
	contentType, bodyReader, err := r.RawBodyStream()
	if err != nil {
		return err
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	switch mediaType {
	case "application/json":
		body, err := io.ReadAll(bodyReader)
		if err != nil {
			return err
		}
		return json.Unmarshal(body, dest)
	case "application/edn":
		body, err := io.ReadAll(bodyReader)
		if err != nil {
			return err
		}
		return edn.Unmarshal(body, dest)
	case "application/x-gzip":
		switch v := dest.(type) {
		case *[]byte:
			gzipReader, err := gzip.NewReader(bodyReader)
			if err != nil {
				return err
			}

			readed := 0
			buffer := make([]byte, 1024)
			for err == nil {
				readed, err = gzipReader.Read(buffer)
				if err != nil && !errors.Is(err, io.EOF) {
					return err
				}
				*v = append(*v, buffer[:readed]...)
			}

			return nil
		default:
			return ErrInvalidUnmarshalTarget
		}
	default:
		return ErrUnmappedMediaType
	}
}

func (r *response) Status() int {
	return r.Response.StatusCode
}

func (r *response) RawBodyStream() (string, io.Reader, error) {
	if r.readed {
		return "", nil, ErrResponseReaded
	}
	r.readed = true
	contentType := r.Response.Header.Get("content-type")
	return contentType, r.Response.Body, nil
}

func (r *response) RawBody() (string, []byte, error) {
	contentType, bodyReader, err := r.RawBodyStream()
	if err != nil {
		return "", nil, err
	}
	body, err := io.ReadAll(bodyReader)
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
