package util

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url string, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dir := filepath.Dir(destination)

	info, err := os.Stat(dir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) || !info.IsDir() {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}

	out, err := os.Create(destination)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	return err
}

func DeleteFile(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}
