package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
	"gorgonia.org/tensor"
)

type sequence struct {
	engine *Engine

	commandPool   vk.CommandPool
	commandBuffer vk.CommandBuffer
	fence         vk.Fence
}

func newSequence(e *Engine) (sequence, error) {
	s := sequence{
		engine: e,
	}
	if err := s.init(); err != nil {
		return sequence{}, err
	}
	return s, nil
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

func (s *sequence) destroy() {
	vk.FreeCommandBuffers(s.engine.device, s.commandPool, 1, []vk.CommandBuffer{s.commandBuffer})
	vk.DestroyCommandPool(s.engine.device, s.commandPool, nil)
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

func (s *sequence) evalSync() error {
	if err := s.evalAsync(true); err != nil {
		return err
	}
	if err := s.await(1000000000); err != nil {
		return err
	}
	return nil
}

func (s *sequence) evalAsync(createFence bool) error {
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

	if createFence {
		fenceInfo := vk.FenceCreateInfo{
			SType: vk.StructureTypeFenceCreateInfo,
		}
		var fence vk.Fence
		res := vk.CreateFence(s.engine.device, &fenceInfo, nil, &fence)
		if res != vk.Success {
			return VulkanError(res)
		}
		s.fence = fence
	}

	res := vk.QueueSubmit(s.engine.computeQueue, 1, []vk.SubmitInfo{submitInfo}, s.fence)
	if res != vk.Success {
		return VulkanError(res)
	}
	return nil
}

// await the eval to be finished, timeout after waitFor nanoseconds
func (s *sequence) await(waitFor uint64) error {
	res := vk.WaitForFences(s.engine.device, 1, []vk.Fence{s.fence}, vk.True, waitFor)
	if res != vk.Success {
		return VulkanError(res)
	}
	vk.DestroyFence(s.engine.device, s.fence, nil)
	s.fence = nil

	return nil
}
