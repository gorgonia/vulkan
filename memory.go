package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
)

//func Allocate(ld *LogicalDevice, bufferSize vk.DeviceSize) error {
//	// TODO: Size must be a multiple of VkPhysicalDeviceLimits.minStorageBufferOffsetAlignment
//	// fmt.Println(ld.physicalDevice.properties.Limits.MinStorageBufferOffsetAlignment)
//	bufferCreateInfo := &vk.BufferCreateInfo{
//		SType:                 vk.StructureTypeBufferCreateInfo,
//		Size:                  bufferSize,
//		Usage:                 vk.BufferUsageFlags(vk.BufferUsageStorageBufferBit),
//		SharingMode:           vk.SharingModeExclusive,
//		QueueFamilyIndexCount: 1,
//		PQueueFamilyIndices:   []uint32{ld.computeQueueFamilyIndex},
//	}
//	var inputBuffer vk.Buffer
//	res := vk.CreateBuffer(ld.device, bufferCreateInfo, nil, &inputBuffer)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	// Check memory requirements
//	var memoryRequirements vk.MemoryRequirements
//	vk.GetBufferMemoryRequirements(ld.device, inputBuffer, &memoryRequirements)
//	memoryRequirements.Deref()
//
//	//fmt.Println(memoryRequirements.Size)
//	//fmt.Println(memoryRequirements.Alignment)
//	//fmt.Println(memoryRequirements.MemoryTypeBits)
//
//	// Allocate memory
//	memoryTypeIndex, err := findMemoryTypeIndex(ld.physicalDevice.device, memoryRequirements, bufferSize)
//	if err != nil {
//		return err
//	}
//
//	var size vk.DeviceSize = 256 * 1024 //memoryRequirements.Size,
//	memoryAllocateInfo := &vk.MemoryAllocateInfo{
//		SType:           vk.StructureTypeMemoryAllocateInfo,
//		AllocationSize:  size,
//		MemoryTypeIndex: memoryTypeIndex,
//	}
//	var memory vk.DeviceMemory
//	res = vk.AllocateMemory(ld.device, memoryAllocateInfo, nil, &memory)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	// Bind buffer
//	res = vk.BindBufferMemory(ld.device, inputBuffer, memory, 0)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	var outputBuffer vk.Buffer
//	res = vk.CreateBuffer(ld.device, bufferCreateInfo, nil, &outputBuffer)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	// Bind buffer
//	res = vk.BindBufferMemory(ld.device, outputBuffer, memory, bufferSize)
//
//	var payload unsafe.Pointer
//	res = vk.MapMemory(ld.device, memory, 0, size, 0, &payload)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	const int32Size = 4
//	// https://stackoverflow.com/a/51188315
//	var payloadSlice []byte
//	sh := (*reflect.SliceHeader)(unsafe.Pointer(&payloadSlice))
//	sh.Data = uintptr(payload)
//	sh.Len = int(size)
//	sh.Cap = int(size)
//	for i := 0; i < int(size)/int32Size; i++ {
//		payloadSlice[i] = byte(rand.Int()) // just for testing purposes
//	}
//
//	vk.UnmapMemory(ld.device, memory)
//
//	// TODO
//
//	code, err := readShaderFile("shaders/compiled/test.spv")
//	if err != nil {
//		return err
//	}
//	shaderModuleCreateInfo := &vk.ShaderModuleCreateInfo{
//		SType:    vk.StructureTypeShaderModuleCreateInfo,
//		CodeSize: uint(len(code)*4), // sizeof(uint32) = 4
//		PCode:    code,
//	}
//	var shaderModule vk.ShaderModule
//	res = vk.CreateShaderModule(ld.device, shaderModuleCreateInfo, nil, &shaderModule)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	descriptorSetLayoutBinding := []vk.DescriptorSetLayoutBinding{
//		{
//			Binding:         0,
//			DescriptorType:  vk.DescriptorTypeStorageBuffer,
//			DescriptorCount: 1,
//			StageFlags:      vk.ShaderStageFlags(vk.ShaderStageComputeBit),
//		},
//		{
//			Binding:         1,
//			DescriptorType:  vk.DescriptorTypeStorageBuffer,
//			DescriptorCount: 1,
//			StageFlags:      vk.ShaderStageFlags(vk.ShaderStageComputeBit),
//		},
//	}
//	descriptorSetLayoutCreateInfo := &vk.DescriptorSetLayoutCreateInfo{
//		SType:        vk.StructureTypeDescriptorSetLayoutCreateInfo,
//		BindingCount: uint32(len(descriptorSetLayoutBinding)),
//		PBindings:    descriptorSetLayoutBinding,
//	}
//	var descriptorSetLayout vk.DescriptorSetLayout
//	res = vk.CreateDescriptorSetLayout(ld.device, descriptorSetLayoutCreateInfo, nil, &descriptorSetLayout)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	pipelineLayoutCreateInfo := &vk.PipelineLayoutCreateInfo{
//		SType:                  vk.StructureTypePipelineLayoutCreateInfo,
//		SetLayoutCount:         1,
//		PSetLayouts:            []vk.DescriptorSetLayout{descriptorSetLayout},
//		PushConstantRangeCount: 0,
//		PPushConstantRanges:    nil,
//	}
//	var pipelineLayout vk.PipelineLayout
//	res = vk.CreatePipelineLayout(ld.device, pipelineLayoutCreateInfo, nil, &pipelineLayout)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	pipelineShaderStageCreateInfo := vk.PipelineShaderStageCreateInfo{
//		SType:               vk.StructureTypePipelineShaderStageCreateInfo,
//		Flags:               0,
//		Stage:               vk.ShaderStageComputeBit,
//		Module:              shaderModule,
//		PName:               "main\x00", // name of the entrypoint function
//		PSpecializationInfo: nil,
//	}
//	computePipelineCreateInfo := []vk.ComputePipelineCreateInfo{
//		{
//			SType:  vk.StructureTypeComputePipelineCreateInfo,
//			Flags:  0,
//			Stage:  pipelineShaderStageCreateInfo,
//			Layout: pipelineLayout,
//		},
//	}
//	pipeline := make([]vk.Pipeline, len(computePipelineCreateInfo))
//	res = vk.CreateComputePipelines(ld.device, vk.NullPipelineCache, uint32(len(computePipelineCreateInfo)), computePipelineCreateInfo, nil, pipeline)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	//
//
//	descriptorPoolSize := vk.DescriptorPoolSize{
//		Type:            vk.DescriptorTypeStorageBuffer,
//		DescriptorCount: 2, // we bind 2 tensors
//	}
//	descriptorPoolCreateInfo := &vk.DescriptorPoolCreateInfo{
//		SType:         vk.StructureTypeDescriptorPoolCreateInfo,
//		Flags:         0,
//		MaxSets:       1,
//		PoolSizeCount: 1,
//		PPoolSizes:    []vk.DescriptorPoolSize{descriptorPoolSize},
//	}
//	var descriptorPool vk.DescriptorPool
//	res = vk.CreateDescriptorPool(ld.device, descriptorPoolCreateInfo, nil, &descriptorPool)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	//
//
//	descriptorSetAllocateInfo := &vk.DescriptorSetAllocateInfo{
//		SType:              vk.StructureTypeDescriptorSetAllocateInfo,
//		DescriptorPool:     descriptorPool,
//		DescriptorSetCount: 1,
//		PSetLayouts:        []vk.DescriptorSetLayout{descriptorSetLayout},
//	}
//	descriptorSet := make([]vk.DescriptorSet, 1)
//	res = vk.AllocateDescriptorSets(ld.device, descriptorSetAllocateInfo, &descriptorSet[0])
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	descriptorBufferInfoInput := vk.DescriptorBufferInfo{
//		Buffer: inputBuffer,
//		Offset: 0,
//		Range:  bufferSize,
//	}
//	descriptorBufferInfoOutput := vk.DescriptorBufferInfo{
//		Buffer: outputBuffer,
//		Offset: 0,
//		Range:  bufferSize,
//	}
//	writeDescriptorSet := vk.WriteDescriptorSet{
//		SType:           vk.StructureTypeWriteDescriptorSet,
//		DstSet:          descriptorSet[0],
//		DstBinding:      0, // starting number of the first element in PBufferInfo. corresponds to "binding" in the shader
//		DstArrayElement: 0,
//		DescriptorCount: 2,
//		DescriptorType:  vk.DescriptorTypeStorageBuffer,
//		PBufferInfo: []vk.DescriptorBufferInfo{
//			descriptorBufferInfoInput,
//			descriptorBufferInfoOutput,
//		},
//	}
//	vk.UpdateDescriptorSets(ld.device, 1, []vk.WriteDescriptorSet{writeDescriptorSet}, 0, nil)
//
//	//
//
//	commandPoolCreateInfo := &vk.CommandPoolCreateInfo{
//		SType:            vk.StructureTypeCommandPoolCreateInfo,
//		Flags:            0,
//		QueueFamilyIndex: ld.computeQueueFamilyIndex,
//	}
//	var commandPool vk.CommandPool
//	res = vk.CreateCommandPool(ld.device, commandPoolCreateInfo, nil, &commandPool)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	commandBufferAllocateInfo := &vk.CommandBufferAllocateInfo{
//		SType:              vk.StructureTypeCommandBufferAllocateInfo,
//		CommandPool:        commandPool,
//		Level:              vk.CommandBufferLevelPrimary,
//		CommandBufferCount: 1,
//	}
//	commandBuffer := make([]vk.CommandBuffer, 1)
//	res = vk.AllocateCommandBuffers(ld.device, commandBufferAllocateInfo, commandBuffer)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	commandBufferBeginInfo := &vk.CommandBufferBeginInfo{
//		SType:            vk.StructureTypeCommandBufferBeginInfo,
//		Flags:            vk.CommandBufferUsageFlags(vk.CommandBufferUsageOneTimeSubmitBit),
//		PInheritanceInfo: nil,
//	}
//	res = vk.BeginCommandBuffer(commandBuffer[0], commandBufferBeginInfo)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	//
//
//	vk.CmdBindPipeline(commandBuffer[0], vk.PipelineBindPointCompute, pipeline[0])
//
//	vk.CmdBindDescriptorSets(commandBuffer[0], vk.PipelineBindPointCompute, pipelineLayout, 0, uint32(len(descriptorSet)), descriptorSet, 0, nil)
//
//	const uint32Size = 4
//	vk.CmdDispatch(commandBuffer[0], uint32(bufferSize) / uint32Size / 512, 1, 1)
//
//	res = vk.EndCommandBuffer(commandBuffer[0])
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	submitInfo := vk.SubmitInfo{
//		SType: vk.StructureTypeSubmitInfo,
//		WaitSemaphoreCount: 0,
//		PWaitSemaphores: nil,
//		PWaitDstStageMask: nil,
//		CommandBufferCount: uint32(len(commandBuffer)),
//		PCommandBuffers: commandBuffer,
//		SignalSemaphoreCount: 0,
//		PSignalSemaphores: nil,
//	}
//	start := time.Now()
//	// TODO: add fence?
//	res = vk.QueueSubmit(ld.computeQueue, 1, []vk.SubmitInfo{submitInfo}, vk.NullFence)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//
//	res = vk.QueueWaitIdle(ld.computeQueue)
//	if res != vk.Success {
//		return VulkanError(res)
//	}
//	end := time.Now()
//	fmt.Println(end.Sub(start))
//
//	//
//
//	res = vk.MapMemory(ld.device, memory, 0, size, 0, &payload)
//
//	start = time.Now()
//	for i := 0; i < int(bufferSize); i++ {
//		if payloadSlice[i] != payloadSlice[int(bufferSize) + i] {
//			return errors.New("buffers do not match :(")
//		}
//	}
//	end = time.Now()
//	fmt.Println(end.Sub(start))
//	fmt.Println("Payloads match! :D")
//
//	vk.UnmapMemory(ld.device, memory)
//
//	//
//
//	vk.FreeCommandBuffers(ld.device, commandPool, 1, commandBuffer)
//	vk.DestroyCommandPool(ld.device, commandPool, nil)
//	//vk.FreeDescriptorSets(ld.device, descriptorPool, 1, &descriptorSet[0]) // needs bit set on pool to allow this
//	vk.DestroyDescriptorPool(ld.device, descriptorPool, nil)
//	vk.DestroyPipeline(ld.device, pipeline[0], nil)
//	vk.DestroyPipelineLayout(ld.device, pipelineLayout, nil)
//	vk.DestroyDescriptorSetLayout(ld.device, descriptorSetLayout, nil)
//	vk.DestroyShaderModule(ld.device, shaderModule, nil)
//	vk.DestroyBuffer(ld.device, inputBuffer, nil)
//	vk.DestroyBuffer(ld.device, outputBuffer, nil)
//	vk.FreeMemory(ld.device, memory, nil)
//
//	return nil
//}

func findMemoryTypeIndex(pd vk.PhysicalDevice, memoryRequirements vk.MemoryRequirements, size vk.DeviceSize) (uint32, error) {
	var memoryProperties vk.PhysicalDeviceMemoryProperties
	vk.GetPhysicalDeviceMemoryProperties(pd, &memoryProperties)
	memoryProperties.Deref()

	for i := range memoryProperties.MemoryHeaps {
		memoryProperties.MemoryHeaps[i].Deref()
	}

	// TODO: change later depending on where the user wants to store their tensor
	memoryPropertyFlags := vk.MemoryPropertyFlags(vk.MemoryPropertyHostVisibleBit | vk.MemoryPropertyHostCoherentBit)

	for i := uint32(0); i < memoryProperties.MemoryTypeCount; i++ {
		if (memoryRequirements.MemoryTypeBits & (1 << i)) == 0 {
			continue
		}

		memoryType := memoryProperties.MemoryTypes[i]
		memoryType.Deref()

		if (memoryType.PropertyFlags&memoryPropertyFlags) == memoryPropertyFlags &&
			memoryProperties.MemoryHeaps[memoryType.HeapIndex].Size >= size {
			// Found our memory type
			return i, nil
		}
	}

	return 0, ErrNoMatchingPhysicalDeviceMemory
}
