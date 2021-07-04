# Vulkan back-end for Gorgonia

[![GoDoc](https://godoc.org/gorgonia.org/vulkan?status.svg)](https://godoc.org/gorgonia.org/vulkan)
[![GitHub version](https://badge.fury.io/gh/gorgonia%2Fvulkan.svg)](https://badge.fury.io/gh/gorgonia%2Fvulkan)
[![Go with Vulkan](https://github.com/gorgonia/vulkan/actions/workflows/go-vulkan.yml/badge.svg)](https://github.com/gorgonia/vulkan/actions/workflows/go-vulkan.yml)
[![codecov](https://codecov.io/gh/gorgonia/vulkan/branch/master/graph/badge.svg)](https://codecov.io/gh/gorgonia/vulkan)
[![Go Report Card](https://goreportcard.com/badge/gorgonia.org/vulkan)](https://goreportcard.com/report/gorgonia.org/vulkan)
[![experimental](http://badges.github.io/stability-badges/dist/experimental.svg)](http://github.com/badges/stability-badges)

This project is in the early development stage.

## Licence
This package can be used under the same license as the main Gorgonia package:<br>
https://github.com/gorgonia/gorgonia#licence

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