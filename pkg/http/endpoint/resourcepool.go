package endpoint

import (
	"compress/gzip"
	"context"

	"github.com/jackc/puddle/v2"
	"github.com/vitorsss/go-helpers/pkg/logs"
)

var gzipReaderPool *puddle.Pool[*gzip.Reader]

func init() {
	var err error
	gzipReaderPool, err = puddle.NewPool(&puddle.Config[*gzip.Reader]{
		Constructor: func(ctx context.Context) (res *gzip.Reader, err error) {
			return &gzip.Reader{}, nil
		},
		Destructor: func(res *gzip.Reader) {
			res.Close()
		},
		MaxSize: 50,
	})
	if err != nil {
		logs.Logger.Error().Err(err).Send()
		panic(err)
	}
}
