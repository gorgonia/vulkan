package vulkan

import (
	"os"
	"reflect"
	"unsafe"
)

func readShaderFile(path string) (spirvData []uint32, err error) {
	var buf []byte
	buf, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(buf)%4 != 0 {
		return nil, ErrSpirvDataNotMultipleOf4Bytes
	}

	hdr := *(*reflect.SliceHeader)(unsafe.Pointer(&buf))
	hdr.Len /= 4
	hdr.Cap /= 4
	spirvData = *(*[]uint32)(unsafe.Pointer(&hdr))

	return
}
