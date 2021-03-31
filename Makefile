all: compile_shaders

compile_shaders:
	$(foreach file, $(wildcard shaders/*.comp), glslc $(file) -fshader-stage=compute -O -o $(file:shaders/%.comp=shaders/compiled/%.spv);)
