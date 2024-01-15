package endpoint

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"

	"github.com/pkg/errors"
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
	ctx      context.Context
	readed   bool
}

func (r *response) Unmarshal(dest interface{}) error {
	contentType, bodyReader, err := r.RawBodyStream()
	if err != nil {
		return err
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return errors.Wrap(err, "failed to parse media type")
	}
	switch mediaType {
	case "application/json":
		body, err := io.ReadAll(bodyReader)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}
		return errors.Wrap(json.Unmarshal(body, dest), "failed to unmarshal json body")
	case "application/edn":
		body, err := io.ReadAll(bodyReader)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}
		return errors.Wrap(edn.Unmarshal(body, dest), "failed to unmarshal edn body")
	case "application/x-gzip":
		switch v := dest.(type) {
		case *[]byte:
			gzipReader, err := gzipReaderPool.Acquire(r.ctx)
			if err != nil {
				return errors.Wrap(err, "failed to acquire gzip reader")
			}
			defer gzipReader.Release()
			err = gzipReader.Value().Reset(bodyReader)
			if err != nil {
				return errors.Wrap(err, "failed to read response body")
			}

			readed := 0
			buffer := make([]byte, 1024)
			for err == nil {
				readed, err = gzipReader.Value().Read(buffer)
				if err != nil && !errors.Is(err, io.EOF) {
					return errors.Wrap(err, "failed to read gziped content")
				}
				*v = append(*v, buffer[:readed]...)
			}

			return nil
		default:
			return errors.Wrap(ErrInvalidUnmarshalTarget, "unmapped")
		}
	default:
		return errors.Wrap(ErrUnmappedMediaType, "unmapped")
	}
}

func (r *response) Status() int {
	return r.Response.StatusCode
}

func (r *response) RawBodyStream() (string, io.Reader, error) {
	if r.readed {
		return "", nil, errors.Wrap(ErrResponseReaded, "")
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
		return "", nil, errors.Wrap(err, "failed to read response body")
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
