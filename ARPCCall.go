package goarpcsolution

import (
	"errors"

	"github.com/AnimusPEXUS/gouuidtools"
)

type ARPCCall struct {
	Name      string
	ReplyToId *gouuidtools.UUID

	Args   []*ARPCCallArg
	CallId *gouuidtools.UUID

	ReplyErrCode uint
	ReplyErrMsg  string
}

func NewARPCCallWithId(
	Name string,
	CallId *gouuidtools.UUID,
	Args ...*ARPCCallArg,
) (*ARPCCall, error) {
	self := new(ARPCCall)
	self.Name = Name
	self.Args = Args
	err := self.IsValidError()
	if err != nil {
		return nil, err
	}
	return self, nil
}

func NewARPCCall(
	Name string,
	Args ...*ARPCCallArg,
) (*ARPCCall, error) {
	return NewARPCCallWithId(Name, nil, Args...)
}

func NewARPCCallReplyWithCodeAndMsg(
	ReplyToId *gouuidtools.UUID,

	ReplyErrCode uint,
	ReplyErrMsg string,

	Args ...*ARPCCallArg,
) (*ARPCCall, error) {
	self := new(ARPCCall)
	self.ReplyToId = ReplyToId
	self.ReplyErrCode = ReplyErrCode
	self.ReplyErrMsg = ReplyErrMsg
	self.Args = Args
	err := self.IsValidError()
	if err != nil {
		return nil, err
	}
	return self, nil
}

// code and msg defaults to 0 and ""
func NewARPCCallReply(
	ReplyToId *gouuidtools.UUID,
	Args ...*ARPCCallArg,
) (*ARPCCall, error) {
	return NewARPCCallReplyWithCodeAndMsg(ReplyToId, 0, "", Args...)
}

func (self *ARPCCall) IsValidError() error {
	if self.Name != "" && self.ReplyToId != nil {
		return errors.New("Name and ReplyToId both present")
	}

	if self.Name == "" && self.ReplyToId == nil {
		return errors.New("Name and ReplyToId both absent")
	}

	return nil
}

func (self *ARPCCall) IsInvalid() bool {
	return self.IsValidError() != nil
}

// ReplyToId != nil
func (self *ARPCCall) HasReplyToIdField() bool {
	return self.ReplyToId != nil
}

// Error != nil
func (self *ARPCCall) HasErrorFields() bool {
	return self.ReplyErrCode != 0 || self.ReplyErrMsg != ""
}

// HasReplyToIdField()
func (self *ARPCCall) HasRequestFields() bool {
	return self.HasReplyToIdField()
}

// self.HasResultField() || self.HasErrorField()
func (self *ARPCCall) HasResponseFields() bool {
	return self.HasReplyToIdField() || self.HasErrorFields()
}

func (self *ARPCCall) IsResponseOrError() bool {

	if self.IsInvalid() {
		return false
	}

	return self.HasResponseFields()
}

func (self *ARPCCall) IsResponseAndNotError() bool {

	if self.IsInvalid() {
		return false
	}

	return self.HasReplyToIdField() && !self.HasErrorFields()
}

// is response and is error
func (self *ARPCCall) IsError() bool {

	if self.IsInvalid() {
		return false
	}

	return !self.HasReplyToIdField() && self.HasErrorFields()
}

func (self *ARPCCall) GenARPCCallForJSON() *ARPCCallForJSON {
	return &ARPCCallForJSON{
		CallId:       self.CallId,
		ReplyToId:    self.ReplyToId,
		Name:         self.Name,
		ReplyErrCode: self.ReplyErrCode,
		ReplyErrMsg:  self.ReplyErrMsg,
	}
}

type ARPCCallForJSON struct {
	CallId       *gouuidtools.UUID
	ReplyToId    *gouuidtools.UUID `json:,omitempty`
	Name         string            `json:,omitempty`
	ReplyErrCode uint              `json:,omitempty`
	ReplyErrMsg  string            `json:,omitempty`
}

type ARPCCallArgSlice struct {
	Args []*ARPCCallArg
}

func (self *ARPCCallArgSlice) NullifyIds() {
	for _, i := range self.Args {
		i.NullifyIds()
	}
}

type ARPCCallArg struct {
	// leave Name empty string if it's a positional arg
	Name string

	Basic           *ARPCCallArgValueTypeBasic
	Buffer          *ARPCCallArgValueTypeBuffer
	Transmission    *ARPCCallArgValueTypeTransmission
	ListeningSocket *ARPCCallArgValueTypeListeningSocket
	ConnectedSocket *ARPCCallArgValueTypeConnectedSocket
}

func (self *ARPCCallArg) IsValidError() error {
	var count uint8 = 0
	for _, i := range []any{
		self.Basic,
		self.Buffer,
		self.Transmission,
		self.ListeningSocket,
		self.ConnectedSocket,
	} {
		if i != nil {
			count++
		}
	}

	if count != 1 {
		return errors.New(
			"must be exactly one of fields of " +
				"ARPCCallArgValue type be defined",
		)
	}

	if self.Buffer != nil && self.Buffer.Id.IsNil() {
		return errors.New("self.Buffer.Id == nil")
	}

	if self.Transmission != nil &&
		self.Transmission.Id.IsNil() {
		return errors.New("self.Transmission.Id == nil")
	}

	if self.ListeningSocket != nil &&
		self.ListeningSocket.Id.IsNil() {
		return errors.New("self.ListeningSocket.Id == nil")
	}

	if self.ConnectedSocket != nil &&
		self.ConnectedSocket.Id.IsNil() {
		return errors.New("self.ConnectedSocket.Id == nil")
	}

	return nil
}

func (self *ARPCCallArg) IsValid() bool {
	return self.IsValidError() == nil
}

func (self *ARPCCallArg) NullifyIds() {
	if self.Buffer != nil {
		self.Buffer.Id = nil
	}
	if self.Transmission != nil {
		self.Transmission.Id = nil
	}
	if self.ListeningSocket != nil {
		self.ListeningSocket.Id = nil
	}
	if self.ConnectedSocket != nil {
		self.ConnectedSocket.Id = nil
	}
}

type ARPCCallArgValueTypeBasic struct {
	OwningArg *ARPCCallArg

	Value any
}

// in following structs, Id used to provide predefined Id and/or to
// return resulting Id, generated and/or used by controller

type ARPCCallArgValueTypeBuffer struct {
	OwningArg *ARPCCallArg

	Id      *gouuidtools.UUID
	Payload ARPCBufferI
}

type ARPCCallArgValueTypeTransmission struct {
	OwningArg *ARPCCallArg

	Id      *gouuidtools.UUID
	Payload ARPCTransmissionI
}

type ARPCCallArgValueTypeListeningSocket struct {
	OwningArg *ARPCCallArg

	Id      *gouuidtools.UUID
	Payload ARPCListeningSocketI
}

type ARPCCallArgValueTypeConnectedSocket struct {
	OwningArg *ARPCCallArg

	Id      *gouuidtools.UUID
	Payload ARPCConnectedSocketI
}

type ARPCCallShortItem struct {
	CallId  *gouuidtools.UUID
	ReplyTo *gouuidtools.UUID
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
	ARPCArgTypeTransmission
	ARPCArgTypeListeningSocket
	ConnectedSocket
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
