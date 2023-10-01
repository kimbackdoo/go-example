package dict

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrNotFound   = errors.New("해당 단어를 찾을 수 없음")
	ErrDuplicated = errors.New("해당 단어가 이미 존재함")
)

type dictionary map[string]string

func CreateDictionary(key string, value string) dictionary {
	return dictionary{key: value}
}

func (dict dictionary) Search(key string) (string, error) {
	value, ok := dict[key]
	if !ok {
		return "", ErrNotFound
	}
	return value, nil
}

func (dict dictionary) Add(key string, value string) error {
	_, err := dict.Search(key)
	switch err {
	case ErrNotFound:
		dict[key] = value
	case nil:
		return ErrDuplicated
	}

	return nil
}

func (dict dictionary) Update(key string, value string) error {
	_, err := dict.Search(key)
	switch err {
	case ErrNotFound:
		return ErrNotFound
	case nil:
		dict[key] = value
	}
	return nil
}

func (dict dictionary) Delete(key string) {
	delete(dict, key)
}

func Example() {
	dict := CreateDictionary("first", "First word")

	value, err := dict.Search("first")
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(value)
	}

	err = dict.Add("second", "Second word")
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(dict.Search("second"))
	}

	err = dict.Update("first", "First word updated")
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(dict.Search("first"))
	}

	dict.Delete("first")
	fmt.Println(dict.Search("first"))
}
