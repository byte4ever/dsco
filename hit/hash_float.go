package hit

import (
	"encoding/binary"
	"hash"
	"reflect"
)

var (
	saltInt   = []byte("int-d4a64668-b28d-4ed2-9760-3f498e663c5e")
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
