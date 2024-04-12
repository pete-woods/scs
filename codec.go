package scs

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Codec is the interface for encoding/decoding session data to and from a byte
// slice for use by the session store.
type Codec interface {
	Encode(created, deadline time.Time, values map[string]interface{}) ([]byte, error)
	Decode([]byte) (created, deadline time.Time, values map[string]interface{}, err error)
}

// GobCodec is used for encoding/decoding session data to and from a byte
// slice using the encoding/gob package.
type GobCodec struct{}

// Encode converts a session deadline and values into a byte slice.
func (GobCodec) Encode(created, deadline time.Time, values map[string]interface{}) ([]byte, error) {
	aux := &struct {
		Created  time.Time
		Deadline time.Time
		Values   map[string]interface{}
	}{
		Created:  created,
		Deadline: deadline,
		Values:   values,
	}

	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(&aux); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Decode converts a byte slice into a session deadline and values.
func (GobCodec) Decode(b []byte) (created, deadline time.Time, values map[string]interface{}, err error) {
	aux := &struct {
		Created  time.Time
		Deadline time.Time
		Values   map[string]interface{}
	}{}

	r := bytes.NewReader(b)
	if err := gob.NewDecoder(r).Decode(&aux); err != nil {
		return time.Time{}, time.Time{}, nil, err
	}

	return aux.Created, aux.Deadline, aux.Values, nil
}
