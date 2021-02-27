package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
	"gorgonia.org/tensor"
)

type sequence struct {
	engine *Engine

	commandPool   vk.CommandPool
	commandBuffer vk.CommandBuffer
}

func (s *sequence) init() error {
	if err := s.createCommandPool(); err != nil {
		return err
	}
	if err := s.createCommandBuffer(); err != nil {
		return err
	}
	return nil
}

func (s *sequence) createCommandPool() error {
	commandPoolInfo := vk.CommandPoolCreateInfo{
		SType:            vk.StructureTypeCommandPoolCreateInfo,
		Flags:            0,
		QueueFamilyIndex: s.engine.computeQueueFamilyIndex,
	}
	var commandPool vk.CommandPool
	res := vk.CreateCommandPool(s.engine.device, &commandPoolInfo, nil, &commandPool)
	if res != vk.Success {
		return VulkanError(res)
	}
	s.commandPool = commandPool
	return nil
}

func (s *sequence) createCommandBuffer() error {
	commandBufferInfo := vk.CommandBufferAllocateInfo{
		SType:              vk.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        s.commandPool,
		Level:              vk.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}
	commandBuffer := make([]vk.CommandBuffer, 1)
	res := vk.AllocateCommandBuffers(s.engine.device, &commandBufferInfo, commandBuffer)
	if res != vk.Success {
		return VulkanError(res)
	}
	s.commandBuffer = commandBuffer[0]
	return nil
}

func (s *sequence) begin() error {
	commandBufferBeginInfo := &vk.CommandBufferBeginInfo{
		SType:            vk.StructureTypeCommandBufferBeginInfo,
		Flags:            vk.CommandBufferUsageFlags(vk.CommandBufferUsageOneTimeSubmitBit),
		PInheritanceInfo: nil,
	}
	res := vk.BeginCommandBuffer(s.commandBuffer, commandBufferBeginInfo)
	return VulkanError(res)
}

func (s *sequence) end() error {
	res := vk.EndCommandBuffer(s.commandBuffer)
	return VulkanError(res)
}

func (s *sequence) record(op Op, params ...tensor.Tensor) error {
	if err := op.Init(params); err != nil {
		return err
	}
	if err := op.Record(); err != nil {
		return err
	}
	return nil
}

func (s *sequence) evalAsync() error {
	submitInfo := vk.SubmitInfo{
		SType:                vk.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   0,
		PWaitSemaphores:      nil,
		PWaitDstStageMask:    nil,
		CommandBufferCount:   1,
		PCommandBuffers:      []vk.CommandBuffer{s.commandBuffer},
		SignalSemaphoreCount: 0,
		PSignalSemaphores:    nil,
	}
	// TODO: add fence
	res := vk.QueueSubmit(s.engine.computeQueue, 1, []vk.SubmitInfo{submitInfo}, vk.NullFence)
	return VulkanError(res)
}
