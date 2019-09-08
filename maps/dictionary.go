package maps

const (
	ErrNotFound   = DictionaryErr("could not find the word you were looking for")
	ErrWordExists = DictionaryErr("could not add the word since it's exists")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
	return string(e)
}

type Dictionary map[string]string

func (d Dictionary) Search(key string) (string, error) {
	def, ok := d[key]
	if !ok {
		return "", ErrNotFound
	}

	return def, nil
}

func (d Dictionary) Add(key, val string) error {
	_, err := d.Search(key)

	switch err {
	case ErrNotFound:
		d[key] = val
	case nil:
		return ErrWordExists
	default:
		return err
	}
	return nil
}
