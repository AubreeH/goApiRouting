package routing

import (
	"fmt"
	"io"
	"os"
	"time"
)

func (f *File) Save() error {
	f.mux.Lock()
	defer f.mux.Unlock()
	return f.save()
}

func (f *File) save() error {
	if f.saved {
		return nil
	}

	err := os.Mkdir("temp", 0777)
	if err != nil {
		return err
	}

	storageLocation := os.Getenv("STORAGE_LOCATION")
	if storageLocation == "" {
		storageLocation = "storage"
	}

	storageLocation += "/temp"

	err = os.MkdirAll(storageLocation, 0777)
	if err != nil {
		return err
	}

	tempFileName := fmt.Sprintf("%s/temp/%d_%s", storageLocation, time.Now().UnixMilli(), f.formFileHeader.Filename)

	file, err := os.Create(tempFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	formFile, err := f.formFileHeader.Open()
	if err != nil {
		return err
	}
	defer formFile.Close()

	_, err = io.Copy(file, formFile)
	if err != nil {
		return err
	}

	f.saved = true

	return nil
}

func (f *File) Store() (string, error) {
	if f.stored {
		return f.FilePath, nil
	}

	if !f.saved {
		err := f.save()
		if err != nil {
			return "", err
		}
	}

	storageLocation := os.Getenv("STORAGE_LOCATION")
	if storageLocation == "" {
		storageLocation = "storage"
	}

	newPath := fmt.Sprintf("%s/%s", storageLocation, f.FileName)

	err := os.Rename(f.FilePath, newPath)
	if err != nil {
		return "", err
	}

	f.FilePath = newPath

	return newPath, nil
}
