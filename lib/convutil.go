package lib

import (
	"bytes"
	"encoding/binary"
	"runtime"
	"strings"
)

func Int16ToBytes(n int16) ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})
	err := binary.Write(buff, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	name := f.Name()

	data := strings.Split(name, ".")
	return data[len(data)-1]
}
