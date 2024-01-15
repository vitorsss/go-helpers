package endpoint

import (
	"compress/gzip"
	"context"

	"github.com/jackc/puddle/v2"
	"github.com/vitorsss/go-helpers/pkg/logs"
)

var gzipPool *puddle.Pool[*gzip.Reader]

func init() {
	var err error
	gzipPool, err = puddle.NewPool(&puddle.Config[*gzip.Reader]{
		Constructor: func(ctx context.Context) (res *gzip.Reader, err error) {
			return nil, nil
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
