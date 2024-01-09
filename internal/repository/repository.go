package repository

import (
	"this_module/internal/repository/filestorage"
)

type Repository struct {
	FileStorage FileStorage
	// Storage
	// Cache
}

type FileStorage interface {
	GetAllFile(string) ([]string, error)
	AddFileInList(string, string)
	DelFileInList(string)
	LockFile(string)
	UnlockFile(string)
	RLockFile(string)
	RUnlockFile(string)
}

func New(filePath string) (*Repository, error) {
	fileStor, err := filestorage.New(filePath)
	if err != nil {
		return nil, err
	}

	return &Repository{
		FileStorage: fileStor,
	}, nil
}

func (r *Repository) GetAllFile(address string) ([]string, error) {
	return r.FileStorage.GetAllFile(address)
}

func (r *Repository) AddFileInList(id, FileExtension string) {
	r.FileStorage.AddFileInList(id, FileExtension)
}

func (r *Repository) DelFileInList(id string) {
	r.FileStorage.DelFileInList(id)
}

func (r *Repository) LockFile(id string) {
	r.FileStorage.LockFile(id)
}

func (r *Repository) UnlockFile(id string) {
	r.FileStorage.UnlockFile(id)
}

func (r *Repository) RLockFile(id string) {
	r.FileStorage.RLockFile(id)
}

func (r *Repository) RUnlockFile(id string) {
	r.FileStorage.RUnlockFile(id)
}
