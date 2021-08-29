package storage

type Storage interface {
	Initialize()
	CreateFileFromURL(name string, url string) error
}
