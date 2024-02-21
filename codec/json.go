package codec

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type JSONCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder
	enc  *json.Encoder
}

func (c *JSONCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *JSONCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *JSONCodec) Close() error {
	return c.conn.Close()
}

func (c *JSONCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()

	err = c.enc.Encode(h)
	if err != nil {
		log.Println("rpc codec: json error encoding header", err)
		return err
	}

	err = c.enc.Encode(body)
	if err != nil {
		log.Println("rpc codec: json error encoding body", err)
		return err
	}

	return nil
}

var _ Codec = (*JSONCodec)(nil)

func NewJSONCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JSONCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf),
	}
}
