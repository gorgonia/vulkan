# This dockerfile is based on the Swiftshader dockerfile contained in the Vulkan Kompute project
# licensed under the Apache License:
# https://github.com/EthicalML/vulkan-kompute/blob/master/docker-builders/Swiftshader.Dockerfile
# https://github.com/EthicalML/vulkan-kompute/blob/master/LICENSE

FROM ubuntu:20.04

# Dependencies
RUN apt-get update -y
RUN apt-get install -y wget gnupg ca-certificates software-properties-common
RUN apt-get install -y git cmake libvulkan-dev

# Dependencies for swiftshader
RUN apt-get install -y g++-8 gcc-8
RUN apt-get install -y libx11-dev zlib1g-dev
RUN apt-get install -y libxext-dev

# Run swiftshader via env VK_ICD_FILENAMES=/swiftshader/vk_swiftshader_icd.json
RUN git clone https://github.com/google/swiftshader swiftshader-build
RUN CC="/usr/bin/gcc-8" CXX="/usr/bin/g++-8" cmake swiftshader-build/. -Bswiftshader-build/build/
# 12 parallel tasks seems to be the magic number big enough that it doesn't
# crash and small enough that it doesn't use all the memory.
RUN cmake --build swiftshader-build/build/. --parallel 12


FROM ubuntu:20.04

# Dependencies
RUN apt-get update -y
RUN apt-get install -y wget curl gcc make gnupg ca-certificates software-properties-common

# Install Node.js so different Go versions can be installed using Github Actions
RUN curl -fsSL https://deb.nodesource.com/setup_16.x | bash -
RUN apt-get install -y nodejs

# Swiftshader
COPY --from=0 swiftshader-build/build/Linux/ /swiftshader/
ENV VK_ICD_FILENAMES=/swiftshader/vk_swiftshader_icd.json

# Install the Vulkan SDK to get access to glslc to compile the shaders
RUN wget -qO - https://packages.lunarg.com/lunarg-signing-key-pub.asc | apt-key add -
RUN wget -qO /etc/apt/sources.list.d/lunarg-vulkan-focal.list https://packages.lunarg.com/vulkan/lunarg-vulkan-focal.list
RUN apt update -y
RUN apt-get install -y vulkan-sdk