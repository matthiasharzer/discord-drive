package storage

import "io"

type Provider interface {
	Save(key string, data io.Reader) error
	Retrieve(key string) (io.ReadCloser, error)
}
