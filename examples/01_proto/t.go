package main

import (
	"fmt"

	"github.com/AnimusPEXUS/gorpcsolution"
)

func main() {
	proto := gorpcsolution.Protocol{
		VersionString: "gorpcsolution protocol v0.0",
		Functions: []gorpcsolution.ProtocolFunction{
			gorpcsolution.ProtocolFunction{
				Name:       "getString",
				Parameters: []gorpcsolution.ProtocolFunctionParameter{},
				ResultType: gorpcsolution.PFRT_String,
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
