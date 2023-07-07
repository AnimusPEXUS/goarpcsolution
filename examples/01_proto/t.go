package main

import (
	"fmt"

	"github.com/AnimusPEXUS/goarpcsolution"
)

func main() {
	proto := gorpcsolution.Protocol{
		VersionString: "gorpcsolution protocol v0.0",
		Functions: []goarpcsolution.ProtocolFunction{
			goarpcsolution.ProtocolFunction{
				Name:       "getString",
				Parameters: []goarpcsolution.ProtocolFunctionParameter{},
				ResultType: gaorpcsolution.PFRT_String,
			},
		},
	}

	m, err := proto.Marshal()
	if err != nil {
		fmt.Println("error marshalling proto:", err)
		return
	}

	fmt.Println("ok:")
	fmt.Println(m)
}
