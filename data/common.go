package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ToJson(p interface{}) ([]byte, error) {

	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func FromJson(data []byte) (Tasks, error) {

	var t Tasks

	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func ReadFile() ([]byte, error) {
	dir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	filepath := dir + "/tasks.json"

	err = checkFile(filepath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}
	return data, nil

}

func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteFile(b []byte) error {

	dir, err := os.Getwd()

	if err != nil {
		return err
	}

	filepath := dir + "/tasks.json"

	tmp := fmt.Sprintf("%s.tmp", filepath)
	unlink(tmp)

	fd, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		return fmt.Errorf("can't create key file %s: %s", tmp, err)
	}

	_, err = fd.Write(b)
	if err != nil {
		fd.Close()
		return fmt.Errorf("can't write %v bytes to %s: %s", len(b), tmp, err)
	}

	fd.Close() // ignore close(2) errors; unrecoverable anyway.

	os.Rename(tmp, filepath)
	return nil
}

func unlink(f string) error {
	st, err := os.Stat(f)
	if err == nil {
		if !st.Mode().IsRegular() {
			return fmt.Errorf("%s can't be unlinked. Not a regular file?", f)
		}

		os.Remove(f)
		return nil
	}

	return err
}
