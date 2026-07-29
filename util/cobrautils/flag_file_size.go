package cobrautils

import (
	"fmt"

	"github.com/docker/go-units"
)

type FlagFileSize struct {
	Bytes int64
}

func (f *FlagFileSize) String() string {
	return units.HumanSize(float64(f.Bytes))
}

func (f *FlagFileSize) Set(value string) error {
	bytes, err := units.RAMInBytes(value)
	if err != nil {
		return fmt.Errorf("invalid file size format: %v", err)
	}
	f.Bytes = bytes
	return nil
}

func (f *FlagFileSize) Type() string {
	return "fileSize"
}
