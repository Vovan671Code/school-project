package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"school-project/data"

	"github.com/djherbis/times"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(page *data.Page) (err error) {

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *data.Page, err error) {

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, data.ErrNoSavedPages
	}

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) PickLast(userName string) (page *data.Page, err error) {

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, data.ErrNoSavedPages
	}
	if len(files) == 1 {
		return s.decodePage(filepath.Join(path, files[0].Name()))
	}

	var site string
	err = os.Chdir(filepath.Join(s.basePath, userName))
	if err != nil {
		fmt.Println(err.Error())
	}
	site = files[0].Name()
	for i, file := range files {
		fmt.Println(i)
		t0, err := times.Stat(file.Name())
		if err != nil {
			fmt.Println("error wirh t0")
		}
		t1, err := times.Stat(site)
		if err != nil {
			fmt.Println("error wirh t1")
		}

		tNow := t0.BirthTime()
		tNext := t1.BirthTime()
		if tNext.Equal(tNow) {
			site = file.Name()
		}
		if tNow.After(tNext) {
			site = file.Name()
		}
	}

	_ = os.Chdir("C:\\Users\\vladi\\Desktop\\school project")

	return s.decodePage(filepath.Join(path, site))
}

func (s Storage) Remove(p *data.Page) error {

	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("can't remove file:%w", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return fmt.Errorf("can't remove file %s, %w", msg, err)
	}
	return nil
}

func (s Storage) IsExists(p *data.Page) (bool, error) {

	fileName, err := fileName(p)
	if err != nil {

		return false, fmt.Errorf("can't check if file exists:%W", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, fmt.Errorf("can't remove file %s, %w", msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*data.Page, error) {
	fmt.Println("decoding page")
	f, err := os.Open(filePath)
	if err != nil {

		return nil, fmt.Errorf("can't decode page; os.open() error: %w", err)
	}
	defer func() { _ = f.Close() }()

	var p data.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("can't decode page with gob decoder: %w", err)
	}

	return &p, nil
}

func fileName(p *data.Page) (string, error) {
	return p.Hash()
}

func (s Storage) ClearAll(userName string) (err error) {
	path := filepath.Join(s.basePath, userName)
	readDirectory, _ := os.Open(path)
	allFiles, _ := readDirectory.Readdir(0)

	for f := range allFiles {
		file := allFiles[f]

		fileName := file.Name()
		filePath := path + "\\" + fileName

		os.Remove(filePath)
	}
	return
}

func (s Storage) PickFirst(userName string) (page *data.Page, err error) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dir)

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, data.ErrNoSavedPages
	}
	if len(files) == 1 {
		return s.decodePage(filepath.Join(path, files[0].Name()))
	}

	var site string
	err = os.Chdir(filepath.Join(s.basePath, userName))
	if err != nil {
		fmt.Println(err.Error())
	}
	site = files[len(files)-1].Name()
	for i, file := range files {
		fmt.Println(i)
		t0, err := times.Stat(file.Name())
		if err != nil {
			fmt.Println("error wirh t0")
		}
		t1, err := times.Stat(site)
		if err != nil {
			fmt.Println("error wirh t1")
		}

		tNow := t0.BirthTime()
		tNext := t1.BirthTime()
		if tNext.Equal(tNow) {
			site = file.Name()
		}
		if tNow.Before(tNext) {
			site = file.Name()
		}
	}

	_ = os.Chdir("C:\\Users\\vladi\\Desktop\\school project")

	return s.decodePage(filepath.Join(path, site))
}
