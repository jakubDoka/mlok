package netw

import (
	"encoding/binary"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestBuffer(t *testing.T) {
	b := Buffer{}

	b.PutBool(false)
	b.PutBool(true)
	b.PutFloat32(1)
	b.PutFloat64(1)
	b.PutInt(1)
	b.PutInt16(1)
	b.PutInt32(1)
	b.PutInt64(1)
	b.PutString("hello")
	b.PutUint(1)
	b.PutUint16(1)
	b.PutUint32(1)
	b.PutUint64(1)

	if b.Bool() || !b.Bool() ||
		b.Float32() != 1 ||
		b.Float64() != 1 ||
		b.Int() != 1 ||
		b.Int16() != 1 ||
		b.Int32() != 1 ||
		b.Int64() != 1 ||
		b.String() != "hello" ||
		b.Uint() != 1 ||
		b.Uint16() != 1 ||
		b.Uint32() != 1 ||
		b.Uint64() != 1 || b.Failed || !b.Finished() {
		b.cursor = 0
		t.Error(
			b.Bool(),
			b.Bool(),
			b.Float32(),
			b.Float64(),
			b.Int(),
			b.Int16(),
			b.Int32(),
			b.Int64(),
			b.String(),
			b.Uint(),
			b.Uint16(),
			b.Uint32(),
			b.Uint64(),
		)
	}
}

func TestInconsistent(t *testing.T) {
	server, client := setup()

	res := Buffer{}
	res.PutString("hello")

	w := Writer{}
	binary.LittleEndian.PutUint32(w.size[:], uint32(len(res.Data)))
	client.Write(w.size[:])

	r := Reader{}
	go func() {
		buff, err := r.Read(server)
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(buff.Data, res.Data) {
			t.Error("\n", buff.Data, "\n", res.Data)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	client.Write(res.Data)

	time.Sleep(time.Millisecond * 100)
}

func TestReadWrite(t *testing.T) {
	server, client := setup()

	res := Buffer{}
	res.PutString("hello")
	go func() {
		var w Writer
		for i := 0; i < 10; i++ {
			_, err := w.Write(client, res.Data)
			if err != nil {
				t.Error(err)
			}
		}
	}()

	var r Reader
	for i := 0; i < 10; i++ {
		buff, err := r.Read(server)
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(buff.Data, res.Data) {
			t.Error("\n", buff.Data, "\n", res.Data)
		}
	}
}

func setup() (c, s net.Conn) {
	adr, err := net.ResolveTCPAddr("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	server, err := net.ListenTCP("tcp", adr)
	if err != nil {
		panic(err)
	}

	adr, err = net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}

	client, err := net.DialTCP("tcp", nil, adr)
	if err != nil {
		panic(err)
	}

	conn, err := server.Accept()
	if err != nil {
		panic(err)
	}

	return client, conn
}
