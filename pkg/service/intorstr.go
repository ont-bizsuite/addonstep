package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"runtime/debug"
	"strconv"
)

type IntOrString struct {
	Type   Type   `protobuf:"varint,1,opt,name=type,casttype=Type"`
	IntVal int64  `protobuf:"varint,2,opt,name=intVal"`
	StrVal string `protobuf:"bytes,3,opt,name=strVal"`
}

// Type represents the stored type of IntOrString.
type Type int

const (
	Int    Type = iota // The IntOrString holds an int.
	String             // The IntOrString holds a string.
)

// FromInt creates an IntOrString object with an int32 value. It is
// your responsibility not to call this method with a value greater
// than int32.
func FromInt(val int64) IntOrString {
	if val > math.MaxInt64 || val < math.MinInt64 {
		log.Printf("value: %d overflows int32\n%s\n", val, debug.Stack())
	}
	return IntOrString{Type: Int, IntVal: int64(val)}
}

// FromString creates an IntOrString object with a string value.
func FromString(val string) IntOrString {
	return IntOrString{Type: String, StrVal: val}
}

// Parse the given string and try to convert it to an integer before
// setting it as a string value.
func Parse(val string) IntOrString {
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return FromString(val)
	}
	return FromInt(i)
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (intstr *IntOrString) UnmarshalJSON(value []byte) error {
	if value[0] == '"' {
		intstr.Type = String
		return json.Unmarshal(value, &intstr.StrVal)
	}
	intstr.Type = Int
	return json.Unmarshal(value, &intstr.IntVal)
}

// String returns the string value, or the Itoa of the int value.
func (intstr *IntOrString) String() string {
	if intstr.Type == String {
		return intstr.StrVal
	}

	return strconv.FormatInt(intstr.IntValue(), 10)
}

// IntValue returns the IntVal if type Int, or if
// it is a String, will attempt a conversion to int.
func (intstr *IntOrString) IntValue() int64 {
	if intstr.Type == String {
		i, _ := strconv.ParseInt(intstr.StrVal, 10, 64)
		return i
	}
	return intstr.IntVal
}

// MarshalJSON implements the json.Marshaller interface.
func (intstr IntOrString) MarshalJSON() ([]byte, error) {
	switch intstr.Type {
	case Int:
		return json.Marshal(intstr.IntVal)
	case String:
		return json.Marshal(intstr.StrVal)
	default:
		return []byte{}, fmt.Errorf("impossible IntOrString.Type")
	}
}
