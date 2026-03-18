package storage

import "io"

type Storage interface {
	Upload(r io.Reader, md5Digest []byte, dstPath string) (string, error)
}
