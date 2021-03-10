package netw

import (
	"reflect"
	"testing"
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

func TestSplitter(t *testing.T) {
	builder := Builder{}

	bs := make([]Buffer, 10)
	for i := range bs {
		b := &bs[i]
		b.PutString("oh no")
		b.PutInt(10)
		b.ops = [8]byte{}
		builder.Write(b.Data)
	}

	spl := Splitter{}
	spl.Recycle(Buffer{})
	spl.Read(builder.Data)

	if !reflect.DeepEqual(bs, spl.buff) {
		for i := range bs {
			if !reflect.DeepEqual(bs[i], spl.buff[i]) {
				t.Errorf("\n%#v\n%#v", bs[i], spl.buff[i])
			}
		}
	}

	spl = Splitter{}
	spl.Read(builder.Data[:20])
	spl.Read(builder.Data[20:])

	if !reflect.DeepEqual(bs, spl.buff) {
		for i := range bs {
			if !reflect.DeepEqual(bs[i], spl.buff[i]) {
				t.Errorf("\n%#v\n%#v", bs[i], spl.buff[i])
			}
		}
	}

	spl = Splitter{}
	spl.Read(builder.Data[:2])
	spl.Read(builder.Data[2:])

	if !reflect.DeepEqual(bs, spl.buff) {
		for i := range bs {
			if !reflect.DeepEqual(bs[i], spl.buff[i]) {
				t.Errorf("\n%#v\n%#v", bs[i], spl.buff[i])
			}
		}
	}

	spl = Splitter{}
	spl.Read(builder.Data[:4])
	spl.Read(builder.Data[4:15])
	spl.Read(builder.Data[15:])

	if !reflect.DeepEqual(bs, spl.buff) {
		for i := range bs {
			if !reflect.DeepEqual(bs[i], spl.buff[i]) {
				t.Errorf("\n%#v\n%#v", bs[i], spl.buff[i])
			}
		}
	}

	builder.Flush()
}

func TestFail(t *testing.T) {
	b := Buffer{}
	if b.Bool() ||
		b.Float32() != 0 ||
		b.Float64() != 0 ||
		b.Int() != 0 ||
		b.Int16() != 0 ||
		b.Int32() != 0 ||
		b.Int64() != 0 ||
		b.String() != "" ||
		b.Uint() != 0 ||
		b.Uint16() != 0 ||
		b.Uint32() != 0 ||
		b.Uint64() != 0 || !b.Failed {
		t.Fail()
	}
}
