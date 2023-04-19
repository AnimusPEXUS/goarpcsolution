package goarpcsolution

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

type ARPCSolutionFuncCallRecord struct {
	Id uuid.UUID
}

type ARPCSolutionFuncArg struct {
	// leave Name empty string if it's a positional arg
	Name  string
	Value ARPCSolutionFuncArgValue
}

type ARPCSolutionFuncRes struct {
}

type ARPCSolutionFuncArgValueTypeBasic struct {
	Value any
}

type ARPCSolutionFuncArgValueTypeBuffer struct {
	Value *ARPCMediaBuffer
}

type ARPCSolutionFuncArgValueTypeBroadcast struct {
	Value *ARPCMediaBroadcast
}

type ARPCSolutionFuncArgValueTypeConn struct {
	Value *ARPCMediaConn
}

type ARPCSolutionFuncArgValue struct {
	Basic     *ARPCSolutionFuncArgValueTypeBasic
	Buffer    *ARPCSolutionFuncArgValueTypeBuffer
	Broadcast *ARPCSolutionFuncArgValueTypeBroadcast
	Conn      *ARPCSolutionFuncArgValueTypeConn
}

func (self *ARPCSolutionFuncArgValue) IsValidError() error {
	var count uint8 = 0
	for _, i := range []any{
		self.Basic,
		self.Buffer,
		self.Broadcast,
		self.Conn,
	} {
		if i != nil {
			count++
		}
	}

	if count != 1 {
		return errors.New("must be exactly one of fields of ARPCSolutionFuncArgValue be defined")
	}

	if self.Buffer != nil && self.Buffer.Value == nil {
		return errors.New("self.Buffer.Value == nil")
	}

	if self.Broadcast != nil && self.Broadcast.Value == nil {
		return errors.New("self.Broadcast.Value == nil")
	}

	if self.Conn != nil && self.Conn.Value == nil {
		return errors.New("self.Conn.Value == nil")
	}

	return nil
}

func (self *ARPCSolutionFuncArgValue) IsValid() bool {
	return self.IsValidError() == nil
}

type ARPCCallShortItem struct {
	CallId  uuid.UUID
	ReplyTo uuid.UUID
}

type ARPCArgType uint8

const (
	ARPCArgTypeInvalid ARPCArgType = iota
	ARPCArgTypeBasicBool
	ARPCArgTypeBasicNumber
	ARPCArgTypeBasicString
	ARPCArgTypeBasicArray
	ARPCArgTypeBasicObject
	ARPCArgTypeBuffer
	ARPCArgTypeBroadcast
	ARPCArgTypePort
)

type ARPCArgInfo struct {
	Type  ARPCArgType
	Value any
}

func (self *ARPCArgInfo) IsValidError() error {
	if self.Type == ARPCArgTypeInvalid {
		return errors.New("invalid type")
	}

	// todo: complete this

	return nil
}

func (self *ARPCArgInfo) IsValid() bool {
	return self.IsValidError() == nil
}
