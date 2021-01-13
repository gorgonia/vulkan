package vulkan

// Some parts of this file were sourced from:
// https://github.com/vulkan-go/asche

import "C"
import (
	vk "github.com/vulkan-go/vulkan"
	"log"
	"unsafe"
)

const layerKhronosValidation = "VK_LAYER_KHRONOS_validation"

func safeString(str string) string {
	const nullStr = "\x00"
	return str + nullStr
}

func safeStrings(strs []string) []string {
	safeStrs := make([]string, len(strs))
	for i, str := range strs {
		safeStrs[i] = safeString(str)
	}
	return safeStrs
}

func debugCallbackFunc(flags vk.DebugReportFlags, objectType vk.DebugReportObjectType,
	object uint64, location uint, messageCode int32, pLayerPrefix string,
	pMessage string, pUserData unsafe.Pointer) vk.Bool32 {

	switch {
	case flags&vk.DebugReportFlags(vk.DebugReportInformationBit) != 0:
		log.Printf("INFORMATION: [%s] Code %d : %s", pLayerPrefix, messageCode, pMessage)
	case flags&vk.DebugReportFlags(vk.DebugReportWarningBit) != 0:
		log.Printf("WARNING: [%s] Code %d : %s", pLayerPrefix, messageCode, pMessage)
	case flags&vk.DebugReportFlags(vk.DebugReportPerformanceWarningBit) != 0:
		log.Printf("PERFORMANCE WARNING: [%s] Code %d : %s", pLayerPrefix, messageCode, pMessage)
	case flags&vk.DebugReportFlags(vk.DebugReportErrorBit) != 0:
		log.Printf("ERROR: [%s] Code %d : %s", pLayerPrefix, messageCode, pMessage)
	case flags&vk.DebugReportFlags(vk.DebugReportDebugBit) != 0:
		log.Printf("DEBUG: [%s] Code %d : %s", pLayerPrefix, messageCode, pMessage)
	default:
		log.Printf("INFORMATION: [%s] Code %d : %s", pLayerPrefix, messageCode, pMessage)
	}
	return vk.Bool32(vk.False)
}
