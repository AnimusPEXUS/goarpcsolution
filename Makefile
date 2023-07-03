export GONOPROXY=github.com/AnimusPEXUS/*

all: get

get: 
		$(MAKE) -C examples/01_proto
		go get -u -v "./..."
		go mod tidy
