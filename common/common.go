package common

import (
	"bytes"
	"encoding/json"
	"net"
)

type Common interface {
	NewDecoder(conn net.Conn) *json.Decoder
	Decode(decoder *json.Decoder) (interface{}, error)
	ReadUntilNewline(conn net.Conn) ([]byte, error)

	Write(conn net.Conn, req []byte) error

	Marshal(v any) ([]byte, error)
}

type commonImpl struct{}

func NewCommonLib() Common {
	return &commonImpl{}
}

func (common *commonImpl) NewDecoder(conn net.Conn) *json.Decoder {
	return json.NewDecoder(conn)
}

func (common *commonImpl) Decode(decoder *json.Decoder) (interface{}, error) {
	var request interface{}
	err := decoder.Decode(&request)

	return request, err
}

func (common *commonImpl) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (common *commonImpl) Write(conn net.Conn, req []byte) error {
	if _, err := conn.Write(append(req, '\n')); err != nil {
		return err
	}

	return nil
}

func (common *commonImpl) ReadUntilNewline(conn net.Conn) ([]byte, error) {
	var buf bytes.Buffer
	for {
		b := make([]byte, 1)
		n, err := conn.Read(b)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			break
		}

		buf.Write(b)
		if b[0] == '\n' {
			break
		}
	}

	return buf.Bytes(), nil
}
