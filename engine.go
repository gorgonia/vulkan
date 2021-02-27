package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
	"gorgonia.org/tensor"
	"unsafe"
)

type Engine struct {
	tensor.Engine

	physicalDevice          vk.PhysicalDevice // Physical device the logical device is on, can be used to get all kind of specs but not much more
	device                  vk.Device         // Logical device, can be used to control the device
	computeQueue            vk.Queue          // Queue to run commands on
	computeQueueFamilyIndex uint32            // Family from which the computeQueue was created

	sequence sequence
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
	if err := e.sequence.evalAsync(); err != nil {
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
	var buffer vk.Buffer
	res := vk.CreateBuffer(e.device, &bufferInfo, nil, &buffer)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	// Find memory requirements
	var requirements vk.MemoryRequirements
	vk.GetBufferMemoryRequirements(e.device, buffer, &requirements)
	requirements.Deref()

	memoryTypeIndex, err := findMemoryTypeIndex(e.physicalDevice, requirements, dSize)
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
	res = vk.BindBufferMemory(e.device, buffer, memory, 0)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	// Map memory so it's accessible
	var pointer unsafe.Pointer
	res = vk.MapMemory(e.device, memory, 0, dSize, 0, &pointer)
	if res != vk.Success {
		return nil, VulkanError(res)
	}

	return &Memory{
		memory:  memory,
		buffer:  buffer,
		pointer: pointer,
		size:    dSize,
	}, nil
}

func (e *Engine) Free(mem tensor.Memory, size int64) error {
	m, ok := mem.(*Memory)
	if !ok {
		return ErrFreeMemoryOfOtherEngine
	}
	if m.pointer == nil {
		return ErrMemoryAlreadyFreed
	}

	vk.UnmapMemory(e.device, m.memory)
	vk.DestroyBuffer(e.device, m.buffer, nil)
	vk.FreeMemory(e.device, m.memory, nil)

	m.memory = vk.NullDeviceMemory
	m.buffer = vk.NullBuffer
	m.pointer = nil
	m.size = 0

	return nil
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

type Memory struct {
	tensor.Memory

	memory  vk.DeviceMemory
	buffer  vk.Buffer
	pointer unsafe.Pointer
	size    vk.DeviceSize
}

func (m *Memory) Uintptr() uintptr {
	return uintptr(m.pointer)
}

func (m *Memory) MemSize() uintptr {
	return uintptr(m.size)
}
