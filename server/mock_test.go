package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

type mockConn struct {
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
	closed      bool
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuffer.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func (m *mockConn) Close() error {
	m.closed = true
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type mockCommon struct {
	shouldReturnErrorOnDecode           bool
	shouldReturnErrorOnMarshal          bool
	shouldReturnErrorOnWrite            bool
	shouldReturnErrorOnReadUntilNewLine bool

	OnNewDecoderCalledCount int
	OnDecodeCalledCount     int
	OnMarshalCalledCount    int
	onWriteCalledCount      int
	onReadUntilNewline      int
}

func (m *mockCommon) NewDecoder(conn net.Conn) *json.Decoder {
	m.OnNewDecoderCalledCount++

	return json.NewDecoder(conn)
}

func (m *mockCommon) Decode(decoder *json.Decoder) (interface{}, error) {
	m.OnDecodeCalledCount++

	if m.shouldReturnErrorOnDecode {
		return nil, io.EOF
	}

	var request interface{}
	err := decoder.Decode(&request)

	return request, err
}

func (m *mockCommon) ReadUntilNewline(conn net.Conn) ([]byte, error) {
	m.onReadUntilNewline++

	if m.shouldReturnErrorOnReadUntilNewLine {
		return nil, io.EOF
	}

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

func (m *mockCommon) Marshal(v any) ([]byte, error) {
	m.OnMarshalCalledCount++

	if m.shouldReturnErrorOnMarshal {
		return nil, fmt.Errorf("mock error")
	}

	return json.Marshal(v)
}

func (m *mockCommon) Write(conn net.Conn, req []byte) error {
	m.onWriteCalledCount++

	if m.shouldReturnErrorOnWrite {
		return fmt.Errorf("mock error")
	}

	if _, err := conn.Write(append(req, '\n')); err != nil {
		return err
	}

	return nil
}
