export GONOPROXY=github.com/AnimusPEXUS/*

all: get

get: 
		make -C examples/01_proto
		go get -u -v "./..."
		go mod tidy
