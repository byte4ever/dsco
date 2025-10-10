package hit

import (
	"encoding/binary"
	"hash"
	"reflect"
)

var (
	saltNil   = []byte("nil-032e68a4-b701-4c25-b0f3-0b33ea9dbcff")
	saltInt   = []byte("int-d4a64668-b28d-4ed2-9760-3f498e663c5e")
	saltInt8  = []byte("int8-d3811dc6-0e84-43e5-b3cf-3d1dd90cf2db")
	saltInt16 = []byte("int16-79f49817-590f-40f8-b91c-346bb93710e6")
	saltInt32 = []byte("int32-1c64ecb1-b276-47d1-a622-d91347afb441")
	saltInt64 = []byte("int64-8c4fb304-d183-44ca-b97b-69b400a3941e")
)

func hdlAnyInt(hasher hash.Hash, value reflect.Value) {
	internalValue := value

	if internalValue.Kind() == reflect.Ptr {
		if internalValue.IsNil() {
			switch internalValue.Type().Elem().Kind() {
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64:
				hasher.Write(saltNil)
				hasher.Write(saltInt)
				hasher.Write(saltNil)

				return

			default:
				panic("not an int")
			}
		}

		internalValue = internalValue.Elem()
	}

	hasher.Write(saltInt)

	buf := make([]byte, binary.MaxVarintLen64)

	binary.PutVarint(buf, internalValue.Int())
	hasher.Write(saltInt)
	hasher.Write(buf)
}
