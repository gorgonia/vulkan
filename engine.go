package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
	"gorgonia.org/tensor"
	"unsafe"
)

type Engine struct {
	logicalDevice
	sequence sequence

	buffers map[unsafe.Pointer]*buffer
}

func NewEngine(pd *Device) (*Engine, error) {
	e := Engine{
		buffers: make(map[unsafe.Pointer]*buffer),
	}
	var err error
	if e.logicalDevice, err = newLogicalDeviceOnPhysicalDevice(pd); err != nil {
		return nil, err
	}
	if e.sequence, err = newSequence(&e); err != nil {
		e.logicalDevice.Destroy()
		return nil, err
	}
	return &e, nil
}

func (e *Engine) Destroy() {
	e.sequence.destroy()
	e.logicalDevice.Destroy()
}

func (e *Engine) evalAsync(op Op, tensors ...tensor.Tensor) error {
	if err := e.sequence.begin(); err != nil {
		return err
	}
	if err := e.sequence.record(op, tensors...); err != nil {
		return err
	}
	if err := e.sequence.end(); err != nil {
		return err
	}
	if err := e.sequence.evalAsync(true); err != nil {
		return err
	}
	return nil
}

func (e *Engine) evalSync(op Op, tensors ...tensor.Tensor) error {
	if err := e.sequence.begin(); err != nil {
		return err
	}
	if err := e.sequence.record(op, tensors...); err != nil {
		return err
	}
	if err := e.sequence.end(); err != nil {
		return err
	}
	if err := e.sequence.evalSync(); err != nil {
		return err
	}
	return nil
}

func (e *Engine) AllocAccessible() bool {
	return true
}

func (e *Engine) Alloc(size int64) (tensor.Memory, error) {
	dSize := vk.DeviceSize(size)

	bufferInfo := vk.BufferCreateInfo{
		SType:                 vk.StructureTypeBufferCreateInfo,
		Size:                  dSize,
		Usage:                 vk.BufferUsageFlags(vk.BufferUsageStorageBufferBit),
		SharingMode:           vk.SharingModeExclusive, // can be accessed from at most 1 queue at once
		QueueFamilyIndexCount: 1,
		PQueueFamilyIndices:   []uint32{e.computeQueueFamilyIndex},
	}
	var buf vk.Buffer
	res := vk.CreateBuffer(e.device, &bufferInfo, nil, &buf)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	// Find memory requirements
	var requirements vk.MemoryRequirements
	vk.GetBufferMemoryRequirements(e.device, buf, &requirements)
	requirements.Deref()

	memoryTypeIndex, err := findMemoryTypeIndex(e.physicalDevice, requirements, dSize) // TODO: change memory type depending on where we want to store the buffer
	if err != nil {
		return nil, err
	}

	// Allocate memory
	memoryInfo := vk.MemoryAllocateInfo{
		SType:           vk.StructureTypeMemoryAllocateInfo,
		AllocationSize:  dSize,
		MemoryTypeIndex: memoryTypeIndex,
	}
	var memory vk.DeviceMemory
	res = vk.AllocateMemory(e.device, &memoryInfo, nil, &memory)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	// Bind buffer to memory
	res = vk.BindBufferMemory(e.device, buf, memory, 0)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	// Map memory so it's accessible
	var pointer unsafe.Pointer
	res = vk.MapMemory(e.device, memory, 0, dSize, 0, &pointer)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	//fmt.Printf("mem: %p buf %p", memory, buffer)

	// Store handles
	handles := &buffer{
		memory: memory,
		buffer: buf,
		size:   dSize,
	}
	//e.buffers[unsafe.Pointer(handles)] = handles // For host-inaccessible memory
	e.buffers[pointer] = handles // For accessible memory

	//return PointerWrapper(unsafe.Pointer(handles)), nil
	return PointerWrapper(pointer), nil
}

func (e *Engine) Free(mem tensor.Memory, size int64) error {
	bufPtr, err := e.handlesFromMemory(mem)
	if err != nil {
		return err
	}
	if bufPtr.size != vk.DeviceSize(size) {
		return ErrPartialMemoryFreeNotSupported
	}

	vk.UnmapMemory(e.device, bufPtr.memory)
	vk.DestroyBuffer(e.device, bufPtr.buffer, nil)
	vk.FreeMemory(e.device, bufPtr.memory, nil)

	delete(e.buffers, unsafe.Pointer(mem.Uintptr()))

	return nil
}

func (e *Engine) FreeTensor(t tensor.Tensor) error {
	mem, err := e.memoryFromTensor(t)
	if err != nil {
		return err
	}
	return e.Free(mem, int64(t.MemSize()))
}

func (e *Engine) Memset(mem tensor.Memory, val interface{}) error {
	panic("not implemented")
}

func (e *Engine) Memclr(mem tensor.Memory) {
	panic("not implemented")
}

func (e *Engine) Memcpy(dst, src tensor.Memory) error {
	// Example:
	// https://github.com/EthicalML/vulkan-kompute/blob/6501c598df112d337cc339e7fca5fcde860234ec/src/Tensor.cpp#L105
	panic("not implemented")
}

func (e *Engine) Accessible(mem tensor.Memory) (tensor.Memory, error) {
	return mem, nil
}

// WorksWith returns false because I haven't looked at this yet
func (e Engine) WorksWith(order tensor.DataOrder) bool {
	return false
}

// NonStdAlloc nothing instead of running the default built in allocator
func (e *Engine) NonStdAlloc() {
}

func (e *Engine) memoryFromTensor(t tensor.Tensor) (tensor.Memory, error) {
	if e != t.Engine() {
		return nil, ErrMemoryManagedByOtherEngine
	}
	return PointerWrapper(t.Uintptr()), nil
}

func (e *Engine) handlesFromTensor(t tensor.Tensor) (*buffer, error) {
	mem, err := e.memoryFromTensor(t)
	if err != nil {
		return nil, err
	}
	return e.handlesFromMemory(mem)
}

func (e *Engine) handlesFromMemory(mem tensor.Memory) (*buffer, error) {
	// Our map contains a valid reference to the buffer so this cast
	// should be ok. The GC may panic if an invalid pointer or previously
	// freed pointer is passed to this function.
	bufPtr := unsafe.Pointer(mem.Uintptr())
	if b, ok := e.buffers[bufPtr]; ok {
		return b, nil
	}
	return nil, ErrUnknownMemory
}

type PointerWrapper uintptr

func (p PointerWrapper) MemSize() uintptr {
	panic("not implemented")
}

// Uintptr returns the wrapped uintptr
func (p PointerWrapper) Uintptr() uintptr {
	return uintptr(p)
}

type buffer struct {
	memory vk.DeviceMemory
	buffer vk.Buffer
	size   vk.DeviceSize
}
