package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
	"gorgonia.org/tensor"
	"unsafe"
)

type spirvData []uint32

const spirvSliceToByteSize = 4 // num of bytes in a uint32
const pushConstantByteSize = 4 // num of bytes in a float32, used for push constants

type Op interface {
	Init(params []tensor.Tensor) error
	Destroy()
	Record() error
}

type opAlgorithmBase struct {
	algorithm
}

func newOpAlgorithmBase(e *Engine) opAlgorithmBase {
	return opAlgorithmBase{
		algorithm: newAlgorithm(e),
	}
}

func (op *opAlgorithmBase) init(shaderFilePath string, params ...tensor.Tensor) error {
	shaderFileData, err := readShaderFile(shaderFilePath)
	if err != nil {
		return err
	}

	return op.algorithm.init(shaderFileData, params)
}

func (op *opAlgorithmBase) destroy() {
	op.algorithm.destroy()
}

type algorithm struct {
	e             *Engine
	pushConstants []int32

	shaderModule        vk.ShaderModule
	descriptorPool      vk.DescriptorPool
	descriptorSetLayout vk.DescriptorSetLayout
	descriptorSet       vk.DescriptorSet

	pipelineLayout vk.PipelineLayout
	pipeline       vk.Pipeline
}

func newAlgorithm(e *Engine) algorithm {
	return algorithm{
		e: e,
	}
}

func (a *algorithm) init(shaderFileData spirvData, params []tensor.Tensor) error {
	if err := a.createParameters(params); err != nil {
		return err
	}
	if err := a.createShaderModule(shaderFileData); err != nil {
		return err
	}
	if err := a.createPipeline(); err != nil {
		return err
	}
	return nil
}

func (a *algorithm) destroy() {
	a.destroyPipeline()
	a.destroyShaderModule()
	a.destroyParameters()
}

func (a *algorithm) createParameters(params []tensor.Tensor) error {
	//e, ok := params[0].Engine().(*Engine)
	//if !ok {
	//	return fmt.Errorf("cannot use tensors that do not belong to the Vulkan engine")
	//}
	//e := params[0].e

	descriptorPoolSize := vk.DescriptorPoolSize{
		Type:            vk.DescriptorTypeStorageBuffer,
		DescriptorCount: uint32(len(params)),
	}
	descriptorPoolInfo := vk.DescriptorPoolCreateInfo{
		SType:         vk.StructureTypeDescriptorPoolCreateInfo,
		Flags:         0,
		MaxSets:       1,
		PoolSizeCount: 1,
		PPoolSizes:    []vk.DescriptorPoolSize{descriptorPoolSize},
	}
	var descriptorPool vk.DescriptorPool
	res := vk.CreateDescriptorPool(a.e.device, &descriptorPoolInfo, nil, &descriptorPool)
	if res != vk.Success {
		return VulkanError(res)
	}
	a.descriptorPool = descriptorPool

	descriptorSetBindings := make([]vk.DescriptorSetLayoutBinding, len(params))
	for i := 0; i < len(params); i++ {
		descriptorSetBindings[i] = vk.DescriptorSetLayoutBinding{
			Binding:         uint32(i),
			DescriptorType:  vk.DescriptorTypeStorageBuffer,
			DescriptorCount: 1,
			StageFlags:      vk.ShaderStageFlags(vk.ShaderStageComputeBit),
		}
	}
	descriptorSetLayoutInfo := vk.DescriptorSetLayoutCreateInfo{
		SType:        vk.StructureTypeDescriptorSetLayoutCreateInfo,
		BindingCount: uint32(len(descriptorSetBindings)),
		PBindings:    descriptorSetBindings,
	}
	var descriptorSetLayout vk.DescriptorSetLayout
	res = vk.CreateDescriptorSetLayout(a.e.device, &descriptorSetLayoutInfo, nil, &descriptorSetLayout)
	if res != vk.Success {
		return VulkanError(res)
	}
	a.descriptorSetLayout = descriptorSetLayout

	descriptorSetInfo := vk.DescriptorSetAllocateInfo{
		SType:              vk.StructureTypeDescriptorSetAllocateInfo,
		DescriptorPool:     descriptorPool,
		DescriptorSetCount: 1,
		PSetLayouts:        []vk.DescriptorSetLayout{descriptorSetLayout},
	}
	var descriptorSet [1]vk.DescriptorSet
	res = vk.AllocateDescriptorSets(a.e.device, &descriptorSetInfo, &descriptorSet[0])
	if res != vk.Success {
		return VulkanError(res)
	}
	a.descriptorSet = descriptorSet[0]

	descriptorBufferInfos := make([]vk.DescriptorBufferInfo, len(params))
	for i, param := range params {
		hndls, err := a.e.handlesFromTensor(param)
		if err != nil {
			return err
		}

		// TODO: move this to Memory
		descriptorBufferInfos[i] = vk.DescriptorBufferInfo{
			Buffer: hndls.buffer,
			Offset: 0,
			Range:  hndls.size,
		}
	}
	writeDescriptorSet := []vk.WriteDescriptorSet{
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          descriptorSet[0],
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: uint32(len(descriptorBufferInfos)),
			DescriptorType:  vk.DescriptorTypeStorageBuffer,
			PBufferInfo:     descriptorBufferInfos,
		},
	}
	vk.UpdateDescriptorSets(a.e.device, 1, writeDescriptorSet, 0, nil)

	return nil
}

func (a *algorithm) destroyParameters() {
	vk.DestroyDescriptorSetLayout(a.e.device, a.descriptorSetLayout, nil)
	vk.DestroyDescriptorPool(a.e.device, a.descriptorPool, nil)
}

func (a *algorithm) createShaderModule(shaderFileData spirvData) error {
	shaderModuleCreateInfo := vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		CodeSize: uint(len(shaderFileData) * spirvSliceToByteSize),
		PCode:    shaderFileData,
	}
	var shaderModule vk.ShaderModule
	res := vk.CreateShaderModule(a.e.device, &shaderModuleCreateInfo, nil, &shaderModule)
	if res != vk.Success {
		return VulkanError(res)
	}
	a.shaderModule = shaderModule
	return nil
}

func (a *algorithm) destroyShaderModule() {
	vk.DestroyShaderModule(a.e.device, a.shaderModule, nil)
}

func (a *algorithm) createPipeline() error {
	pipelineLayoutInfo := vk.PipelineLayoutCreateInfo{
		SType:                  vk.StructureTypePipelineLayoutCreateInfo,
		SetLayoutCount:         1,
		PSetLayouts:            []vk.DescriptorSetLayout{a.descriptorSetLayout},
		PushConstantRangeCount: 0,
		PPushConstantRanges:    nil,
	}

	if len(a.pushConstants) > 0 {
		pipelineLayoutInfo.PushConstantRangeCount = 1
		pipelineLayoutInfo.PPushConstantRanges = []vk.PushConstantRange{
			{
				StageFlags: vk.ShaderStageFlags(vk.ShaderStageComputeBit),
				Offset:     0,
				Size:       uint32(pushConstantByteSize * len(a.pushConstants)),
			},
		}
	}

	var pipelineLayout vk.PipelineLayout
	res := vk.CreatePipelineLayout(a.e.device, &pipelineLayoutInfo, nil, &pipelineLayout)
	if res != vk.Success {
		return VulkanError(res)
	}

	// TODO: specialization info

	pipelineShaderStageInfo := vk.PipelineShaderStageCreateInfo{
		SType:               vk.StructureTypePipelineShaderStageCreateInfo,
		Flags:               0,
		Stage:               vk.ShaderStageComputeBit,
		Module:              a.shaderModule,
		PName:               "main\x00", // null terminated name of the entrypoint function
		PSpecializationInfo: nil,
	}
	computePipelineInfo := []vk.ComputePipelineCreateInfo{
		{
			SType:  vk.StructureTypeComputePipelineCreateInfo,
			Flags:  0,
			Stage:  pipelineShaderStageInfo,
			Layout: pipelineLayout,
		},
	}
	// TODO: pipeline cache?
	pipeline := make([]vk.Pipeline, len(computePipelineInfo))
	res = vk.CreateComputePipelines(a.e.device, vk.NullPipelineCache, uint32(len(computePipelineInfo)), computePipelineInfo, nil, pipeline)
	if res != vk.Success {
		return VulkanError(res)
	}
	a.pipelineLayout = pipelineLayout
	a.pipeline = pipeline[0]

	return nil
}

func (a *algorithm) destroyPipeline() {
	vk.DestroyPipeline(a.e.device, a.pipeline, nil)
	vk.DestroyPipelineLayout(a.e.device, a.pipelineLayout, nil)
}

func (a *algorithm) recordDispatch(x uint32, y uint32, z uint32) {
	// recordBindCore
	vk.CmdBindPipeline(a.e.sequence.commandBuffer, vk.PipelineBindPointCompute, a.pipeline)
	vk.CmdBindDescriptorSets(a.e.sequence.commandBuffer, vk.PipelineBindPointCompute, a.pipelineLayout, 0, 1, []vk.DescriptorSet{a.descriptorSet}, 0, nil)

	// recordBindPush
	if len(a.pushConstants) > 0 {
		vk.CmdPushConstants(
			a.e.sequence.commandBuffer,
			a.pipelineLayout,
			vk.ShaderStageFlags(vk.ShaderStageComputeBit),
			0,
			uint32(pushConstantByteSize * len(a.pushConstants)),
			unsafe.Pointer(&a.pushConstants[0]),
		)
	}

	// recordDispatch
	vk.CmdDispatch(a.e.sequence.commandBuffer, x, y, z)
}
