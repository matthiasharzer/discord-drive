package storage

import "io"

type Provider interface {
	Write(key string, data io.Reader) error
	Read(key string) (io.ReadCloser, error)
}
