package config

import (
	"bytes"
	"strings"
)

type CSV []string

//goland:noinspection ALL
func (pk *CSV) UnmarshalText(text []byte) error {
	bs := bytes.Split(text, []byte(","))
	*pk = make([]string, 0)
	for _, b := range bs {
		*pk = append(*pk, string(bytes.TrimSpace(b)))
	}
	return nil
}

//goland:noinspection ALL
func (pk CSV) String() string {
	return strings.Join(pk, ", ")
}
