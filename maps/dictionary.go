package maps

import "errors"

var ErrNotFound = errors.New("could not find the word you were looking for")

type Dictionary map[string]string

func (d Dictionary) Search(key string) (string, error) {
	def, ok := d[key]
	if !ok {
		return "", ErrNotFound
	}

	return def, nil
}

func (d Dictionary) Add(key, definition string) {
	d[key] = definition
}
