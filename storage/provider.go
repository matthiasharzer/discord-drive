package storage

import "io"

type Provider interface {
	io.Closer
	Write(key string, data io.Reader) error
	Read(key string) (io.ReadCloser, error)
	Has(key string) (bool, error)
}
