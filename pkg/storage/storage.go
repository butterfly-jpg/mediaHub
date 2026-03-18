package storage

type Storage interface {
	Upload(content string) error
}
