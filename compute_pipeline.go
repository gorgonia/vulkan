package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
	"os"
	"unsafe"
)

func readShaderFile(path string) (spirvData []uint32, err error) {
	var buf []byte
	buf, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(buf) % 4 != 0 {
		return nil, ErrSpirvDataNotMultipleOf4Bytes
	}

	spirvData = make([]uint32, len(buf)/4)
	vk.Memcopy(unsafe.Pointer(&spirvData[0]), buf)
	return
}
