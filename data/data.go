package data

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
	PickFirst(userName string) (*Page, error)
	PickLast(userName string) (*Page, error)
	ClearAll(userName string) error
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", fmt.Errorf("can't calculate hash:%w", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", fmt.Errorf("can't calculate hash:%w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
