package common

import (
	"encoding/json"
	"net"
)

type Common interface {
	NewDecoder(conn net.Conn) *json.Decoder
	Decode(decoder *json.Decoder) (interface{}, error)
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
