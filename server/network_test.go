package server

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockNetwork struct {
	onHandleConnectionCount int
}

func (m *mockNetwork) HandleConnection(ctx context.Context, conn net.Conn, callback func(ctx context.Context, req []byte) interface{}) {
	m.onHandleConnectionCount++
}

func TestHandleConnection_SUCCESS(t *testing.T) {
	t.Parallel()
	mockLib := &mockCommon{}

	newNetwork := &networkImpl{
		common: mockLib,
	}

	conn := &mockConn{
		readBuffer:  bytes.NewBufferString("{}"),
		writeBuffer: &bytes.Buffer{},
	}

	callbackWasCalled := false
	callback := func(ctx context.Context, req []byte) interface{} {
		callbackWasCalled = true
		return "result-mock"
	}

	ctx, cancel := context.WithCancel(context.Background())
	go newNetwork.HandleConnection(ctx, conn, callback)

	time.Sleep(2 * time.Second)
	cancel()

	assert.GreaterOrEqual(t, mockLib.OnNewDecoderCalledCount, 1, "Expected New Decoder function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnDecodeCalledCount, 1, "Expected Decode function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnMarshalCalledCount, 1, "Expected Marshal function to be called at least one time")
	assert.Equal(t, callbackWasCalled, true, "Expected callback to be called at least 1 time")
	assert.NotEmpty(t, conn.writeBuffer.String(), "The connection result should be empty")
	assert.Contains(t, conn.writeBuffer.String(), "result-mock", "The connection result should be the same as the callback result")
}

func TestHandleConnection_ERROR_Context_Closed(t *testing.T) {
	t.Parallel()
	mockLib := &mockCommon{}
	newNetwork := &networkImpl{
		common: mockLib,
	}

	conn := &mockConn{
		readBuffer:  bytes.NewBufferString("some data"), //invalid json data
		writeBuffer: &bytes.Buffer{},
	}

	callbackWasCalled := false
	callback := func(ctx context.Context, req []byte) interface{} {
		callbackWasCalled = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	go newNetwork.HandleConnection(ctx, conn, callback)

	time.Sleep(1 * time.Second)

	assert.GreaterOrEqual(t, mockLib.OnNewDecoderCalledCount, 1, "Expected New Decoder function to be called at least one time")
	assert.Equal(t, mockLib.OnDecodeCalledCount, 0, "Expected Decode function to not be called at all")
	assert.Equal(t, callbackWasCalled, false, "Expected callback to be called 0 times")
}

func TestHandleConnection_ERROR_Invalid_Json(t *testing.T) {
	t.Parallel()
	mockLib := &mockCommon{}
	newNetwork := &networkImpl{
		common: mockLib,
	}

	conn := &mockConn{
		readBuffer:  bytes.NewBufferString("some data"), //invalid json data
		writeBuffer: &bytes.Buffer{},
	}

	callbackWasCalled := false
	callback := func(ctx context.Context, req []byte) interface{} {
		callbackWasCalled = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go newNetwork.HandleConnection(ctx, conn, callback)

	time.Sleep(1 * time.Second)
	defer cancel()

	assert.GreaterOrEqual(t, mockLib.OnNewDecoderCalledCount, 1, "Expected New Decoder function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnDecodeCalledCount, 1, "Expected Decode function to be called at least one time")
	assert.Equal(t, callbackWasCalled, false, "Expected callback to be called 0 times")
}

func TestHandleConnection_ERROR_Lost_Connection_IOF(t *testing.T) {
	t.Parallel()
	mockLib := &mockCommon{
		shouldReturnErrorOnDecode: true,
	}

	newNetwork := &networkImpl{
		common: mockLib,
	}

	conn := &mockConn{
		readBuffer:  bytes.NewBufferString("{}"),
		writeBuffer: &bytes.Buffer{},
	}

	callbackWasCalled := false
	callback := func(ctx context.Context, req []byte) interface{} {
		callbackWasCalled = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go newNetwork.HandleConnection(ctx, conn, callback)

	time.Sleep(1 * time.Second)
	defer cancel()

	assert.GreaterOrEqual(t, mockLib.OnNewDecoderCalledCount, 1, "Expected New Decoder function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnDecodeCalledCount, 1, "Expected Decode function to be called at least one time")
	assert.Equal(t, callbackWasCalled, false, "Expected callback to be called 0 times")
}

func TestHandleConnection_ERROR_Marshal(t *testing.T) {
	t.Parallel()
	mockLib := &mockCommon{
		shouldReturnErrorOnMarshal: true,
	}

	newNetwork := &networkImpl{
		common: mockLib,
	}

	conn := &mockConn{
		readBuffer:  bytes.NewBufferString("{}"),
		writeBuffer: &bytes.Buffer{},
	}

	callbackWasCalled := false
	callback := func(ctx context.Context, req []byte) interface{} {
		callbackWasCalled = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go newNetwork.HandleConnection(ctx, conn, callback)

	time.Sleep(1 * time.Second)
	defer cancel()

	assert.GreaterOrEqual(t, mockLib.OnNewDecoderCalledCount, 1, "Expected New Decoder function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnDecodeCalledCount, 1, "Expected Decode function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnMarshalCalledCount, 1, "Expected Marshal function to be called at least one time")
	assert.Equal(t, callbackWasCalled, false, "Expected callback to be called 0 times")
	assert.Equal(t, conn.writeBuffer.String(), "", "The connection result should be empty")
}

func TestHandleConnection_ERROR_Marshal_Result(t *testing.T) {
	t.Parallel()
	mockLib := &mockCommon{}

	newNetwork := &networkImpl{
		common: mockLib,
	}

	conn := &mockConn{
		readBuffer:  bytes.NewBufferString("{}"),
		writeBuffer: &bytes.Buffer{},
	}

	callbackWasCalled := false
	callback := func(ctx context.Context, req []byte) interface{} {
		callbackWasCalled = true

		type Invalid struct {
			Ch chan int
		}

		data := Invalid{
			Ch: make(chan int),
		}

		return data
	}

	ctx, cancel := context.WithCancel(context.Background())
	go newNetwork.HandleConnection(ctx, conn, callback)

	time.Sleep(2 * time.Second)
	cancel()

	assert.GreaterOrEqual(t, mockLib.OnNewDecoderCalledCount, 1, "Expected New Decoder function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnDecodeCalledCount, 1, "Expected Decode function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnMarshalCalledCount, 1, "Expected Marshal function to be called at least one time")
	assert.Equal(t, callbackWasCalled, true, "Expected callback to be called at least 1 time")
	assert.Empty(t, conn.writeBuffer.String(), "The connection result should be empty")
}

func TestHandleConnection_ERROR_WriteConn(t *testing.T) {
	t.Parallel()
	mockLib := &mockCommon{
		shouldReturnErrorOnWrite: true,
	}

	newNetwork := &networkImpl{
		common: mockLib,
	}

	conn := &mockConn{
		readBuffer:  bytes.NewBufferString("{}"),
		writeBuffer: &bytes.Buffer{},
	}

	callbackWasCalled := false
	callback := func(ctx context.Context, req []byte) interface{} {
		callbackWasCalled = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go newNetwork.HandleConnection(ctx, conn, callback)

	time.Sleep(1 * time.Second)
	defer cancel()

	assert.GreaterOrEqual(t, mockLib.OnNewDecoderCalledCount, 1, "Expected New Decoder function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnDecodeCalledCount, 1, "Expected Decode function to be called at least one time")
	assert.GreaterOrEqual(t, mockLib.OnMarshalCalledCount, 1, "Expected Marshal function to be called at least one time")
	assert.Equal(t, callbackWasCalled, true, "Expected callback to be called at least 1 time")
	assert.Equal(t, conn.writeBuffer.String(), "", "The connection result should be empty")
}
