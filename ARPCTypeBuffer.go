package goarpcsolution

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AnimusPEXUS/gouuidtools"
)

type ARPCBufferMode uint8

const (
	ARPCBufferModeInvalid ARPCBufferMode = iota

	// if buffer items are byte slices and passed throug calls
	// as byte arrays
	ARPCBufferModeBinary

	// if buffer items are objects and passed through calls
	// as JS objects
	ARPCBufferModeObject
)

type ARPCBufferInfo struct {
	Id               *gouuidtools.UUID
	HumanTitle       string
	HumanDescription string

	// todo: todo. this part is in fog for now
	Mode            ARPCBufferMode
	Finished        bool
	TechDescription any
}

type ARPCBufferItem struct {
	BufferId *gouuidtools.UUID
	ItemId   string
	ItemTime time.Time

	Value any
}

type ARPCBufferI interface {
	GetInfo() *ARPCBufferInfo
	ItemCount() int
	// if not found - it's not error and 2nd result is false
	GetItem(id string) (*ARPCBufferItem, bool, error)
}

type ARPCBufferItemSpecifierType uint8

const (
	ARPCBufferItemSpecifierTypeInvalid ARPCBufferItemSpecifierType = iota
	// ARPCBufferItemSpecifierTypeUUID
	ARPCBufferItemSpecifierTypeIndex
	ARPCBufferItemSpecifierTypeTime
	ARPCBufferItemSpecifierTypeString
)

type ARPCBufferItemSpecifier struct {
	Value string
}

func NewARPCBufferItemSpecifierFromString(
	value string,
) (*ARPCBufferItemSpecifier, ARPCBufferItemSpecifierType) {
	self := new(ARPCBufferItemSpecifier)
	self.Value = strings.TrimSpace(value)
	t, _ := self.Type()
	return self, t
}

// var RE_UUID_Id = regexp.MustCompile(`^[0-9a-fA-F_\-]{32,}$`)
// var RE_Index = regexp.MustCompile(`^\-?\d+$`)

func (self *ARPCBufferItemSpecifier) Type() (
	ARPCBufferItemSpecifierType,
	string,
) {

	var ret ARPCBufferItemSpecifierType = ARPCBufferItemSpecifierTypeInvalid
	var xvals []string

	xval := self.Value

	xval = strings.TrimSpace(xval)

	len_xval := len(xval)

	if len_xval == 0 {
		goto ret_invalid
	}

	xvals = strings.SplitN(xval, ":", 2)

	if len(xvals) != 2 {
		goto ret_invalid
	}

	switch xvals[0] {
	default:
		goto ret_invalid
	// case "U":
	// 	ret = ARPCBufferItemSpecifierTypeUUID
	case "#":
		ret = ARPCBufferItemSpecifierTypeIndex
	case "T":
		ret = ARPCBufferItemSpecifierTypeTime
	case "S":
		ret = ARPCBufferItemSpecifierTypeString
	}

	return ret, xvals[1]

ret_invalid:

	return ARPCBufferItemSpecifierTypeInvalid, ""
}

func (self *ARPCBufferItemSpecifier) Index() (int, bool) {

	if t, s := self.Type(); t == ARPCBufferItemSpecifierTypeIndex {
		res_int, err := strconv.Atoi(s)
		if err == nil {
			return res_int, true
		}
	}

	return 0, false
}

func (self *ARPCBufferItemSpecifier) SetIndex(v int) {
	self.Value = fmt.Sprintf("#:%d", v)
}

func (self *ARPCBufferItemSpecifier) Time() (time.Time, bool) {

	if t, s := self.Type(); t == ARPCBufferItemSpecifierTypeTime {
		x, err := time.Parse(time.RFC3339Nano, s)
		if err == nil {
			return x, true
		}
	}

	return time.Time{}, false
}

func (self *ARPCBufferItemSpecifier) SetTime(t time.Time) {
	self.Value = fmt.Sprintf("T:%s", t.Format(time.RFC3339Nano))
}

// func (self *ARPCBufferItemSpecifier) UUID() (uuid.UUID, bool) {

// 	if t, s := self.Type(); t == ARPCBufferItemSpecifierTypeUUID {
// 		for _, s1 := range []string{"-", "_", " "} {
// 			for strings.Contains(s, s1) {
// 				s = strings.Replace(s, s1, "", -1)
// 			}
// 		}
// 		h, err := hex.DecodeString(s)
// 		if err == nil {
// 			return uuid.FromBytesOrNil(h), true
// 		}
// 	}

// 	return uuid.Nil, false
// }

// func (self *ARPCBufferItemSpecifier) SetUUID(v uuid.UUID) {
// 	self.Value = fmt.Sprintf("U:%s", strings.ToLower(v.String()))
// }

func (self *ARPCBufferItemSpecifier) StringVal() (string, bool) {

	if t, s := self.Type(); t == ARPCBufferItemSpecifierTypeString {
		return s, true
	}

	return "", false
}

func (self *ARPCBufferItemSpecifier) SetStringVal(s string) {
	self.Value = fmt.Sprintf("S:%s", s)
}
