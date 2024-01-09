package filestorage

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FStorage struct {
	Files map[string]*File
}

type File struct {
	mu   sync.RWMutex
	Name string
}

func New(filePath string) (*FStorage, error) {
	s := &FStorage{
		Files: make(map[string]*File),
	}

	files, err := s.GetAllFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(files) > 0 {
		for _, v := range files {
			name := strings.Split(v, ".")
			s.Files[name[0]] = &File{Name: v}
		}
	}

	return s, nil
}

func (s *FStorage) GetAllFile(fileAddress string) ([]string, error) {
	var files []string

	err := filepath.Walk(fileAddress, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (s *FStorage) AddFileInList(name, extension string) {
	s.Files[name] = &File{Name: name + extension}
}

func (s *FStorage) DelFileInList(name string) {
	delete(s.Files, name)
}

func (s *FStorage) LockFile(id string) {
	file, ok := s.Files[id]
	if ok {
		file.mu.Lock()
	}
}

func (s *FStorage) UnlockFile(id string) {
	file, ok := s.Files[id]
	if ok {
		file.mu.Unlock()
	}
}

func (s *FStorage) RLockFile(id string) {
	file, ok := s.Files[id]
	if ok {
		file.mu.RLock()
	}
}

func (s *FStorage) RUnlockFile(id string) {
	file, ok := s.Files[id]
	if ok {
		file.mu.RUnlock()
	}
}
