package go8583

import (
	"errors"
	"strconv"

	"github.com/doswell/go8583/util"
)

type variableField struct {
	varLength int
}

//Unpack the field from the message
func (f *variableField) Unpack(offset int, data []byte) (newOffset int, fieldData []byte, err error) {
	//Read varLength to get length.,
	fieldEnd := offset + f.varLength

	if fieldEnd > len(data) {
		return fieldEnd, fieldData, errors.New("Attempt to read passed end of data")
	}
	//Assumes lengths are packed as ASCII. Lengths could be packed as bytes depending on ISO8583 formatting.
	fieldSize, err := strconv.Atoi(string(data[offset:fieldEnd]))
	fieldStart := fieldEnd
	fieldEnd = fieldEnd + fieldSize
	if fieldEnd > len(data) {
		return 0, fieldData, errors.New("Attempt to read passed end of data")
	}
	fieldData = data[fieldStart:fieldEnd]
	//value = FieldValue{Value: string(fieldData)}
	return fieldEnd, fieldData, nil
}

func (f *variableField) Pack(fieldData []byte) (data []byte, err error) {
	//	if isOk, err := f.ValidateValue(value.Value); !isOk {
	//		return nil, err
	//	}

	msgData := []byte(util.LeftPad2Len(strconv.Itoa(len(fieldData)), "0", f.varLength))
	msgData = append(msgData, fieldData...)
	return msgData, nil
}

//NewLVarField creates a new variable length field denoted by a length of one byte. Valid lengths are 0-9
func NewLVarField(bitNumber int, name string, size int, fieldType fieldType) Field {
	return &BitmapMessageField{bitNumber, name, fieldType, LVar, size, NewVariableFieldPackerUnpacker(1), nil}
}

//NewLlVarField creates a new variable length field denoted by a length of two bytes. Valid lengths are 00-99
func NewLlVarField(bitNumber int, name string, size int, fieldType fieldType) Field {
	return &BitmapMessageField{bitNumber, name, fieldType, LlVar, size, NewVariableFieldPackerUnpacker(2), nil}
}

//NewLllVarField creates a new variable length field denoted by a length of three bytes. Valid lengths are 000-999
func NewLllVarField(bitNumber int, name string, size int, fieldType fieldType) Field {
	return &BitmapMessageField{bitNumber, name, fieldType, LlVar, size, NewVariableFieldPackerUnpacker(3), nil}
}

//NewLlllVarField creates a new variable length field denoted by a length of four bytes. Valid lengths are 0000-9999
func NewLlllVarField(bitNumber int, name string, size int, fieldType fieldType) Field {
	return &BitmapMessageField{bitNumber, name, fieldType, LlVar, size, NewVariableFieldPackerUnpacker(4), nil}
}

//NewLllllVarField creates a new variable length field denoted by a length of five bytes. Valid lengths are 00000-99999
func NewLllllVarField(bitNumber int, name string, size int, fieldType fieldType) Field {
	return &BitmapMessageField{bitNumber, name, fieldType, LlVar, size, NewVariableFieldPackerUnpacker(5), nil}
}

//NewLlllllVarField creates a new variable length field denoted by a length of six bytes. Valid lengths are 000000-999999
func NewLlllllVarField(bitNumber int, name string, size int, fieldType fieldType) Field {
	return &BitmapMessageField{bitNumber, name, fieldType, LlVar, size, NewVariableFieldPackerUnpacker(6), nil}
}

//NewLllllVarFieldWithCustomUnpacker creates a new variable length field denoted by a length of five bytes. Valid lengths are 000000-999999, and using a custom unpacker.
func NewLllllVarFieldWithCustomUnpacker(bitNumber int, name string, size int, fieldType fieldType, fieldPackerUnpacker FieldPackerUnpacker) Field {
	f := &BitmapMessageField{bitNumber, name, fieldType, LlVar, size, NewVariableFieldPackerUnpacker(6), fieldPackerUnpacker}
	f.FieldPackerUnpacker = fieldPackerUnpacker
	return f
}

//NewVariableFieldPackerUnpacker creates a new variable field
func NewVariableFieldPackerUnpacker(size int) PackerUnpacker {

	return &variableField{varLength: size}
	// return func(offset int, data []byte) (newOffset int, fieldData []byte, err error) {
	// 		newOffset = offset + size
	// 		if newOffset > len(data) {
	// 			return newOffset, nil, errors.New("Exceeded bytes")
	// 		}

	// 		dataLength, err := strconv.Atoi(string(data[offset:newOffset]))
	// 		if err != nil {
	// 			return newOffset, nil, err
	// 		}
	// 		fieldDataStart := newOffset
	// 		newOffset = newOffset + dataLength
	// 		if newOffset > len(data) {
	// 			return newOffset, nil, errors.New("Exceeded bytes")
	// 		}

	// 		return newOffset, data[fieldDataStart:newOffset], nil
	// 	}, func(data []byte) (msgData []byte, err error) {
	// 		return nil
	// 	}
}
