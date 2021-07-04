# Vulkan back-end for Gorgonia

This project is in the early development stage.

## Dependencies
|Package|Used For|Vitality|Notes|Licence|
|-------|--------|--------|-----|-------|
|[vulkan-go/vulkan](https://github.com/vulkan-go/vulkan) | Making calls to Vulkan | Vital | | [MIT](https://github.com/vulkan-go/vulkan/blob/master/LICENSE.txt) |
|[google/swiftshader](https://github.com/google/swiftshader) | Testing Vulkan in Github Actions without GPU | Only used for testing. A computer or server with GPU can also be used | |[Apache-2.0](https://github.com/google/swiftshader/blob/master/LICENSE.txt) |

## Various Other Copyright Notices
These are the packages and libraries which inspired and were adapted from
in the process of writing this package:

| Source | How it's Used | Licence |
|------|---|-------|
| Vulkan Kompute  | Used as a reference for Vulkan code.<br>The Swiftshader docker image used in the CI pipeline is based on one of its docker images. | [Apache-2.0](https://github.com/EthicalML/vulkan-kompute/blob/master/LICENSE) |