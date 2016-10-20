package go8583

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/doswell/go8583/util"
)

type MessageTemplate interface {
	GetFieldDef(fieldNumber int) (field Field, err error)
}

type Message interface {
	SetField(fieldNr int, value FieldValue)
	SetString(fieldNr int, value string)
	GetField(fieldNr int) (string, bool)
	GetMsgType() int
	Pack() ([]byte, error)
	SetMsgType(msgType int)
	CopyField(fieldNr int, msg Message)
	GetSubField(fieldNr int, subFieldNr int) (string, bool)
	SetSubField(fieldNr, subFieldNr int, subValue string)
	GetMsgTypeString() string
}

type BitmapMessageTemplate struct {
	Header []Field
	Fields map[int]Field
	lock   sync.Mutex
}

type BitmapMessage struct {
	*BitmapMessageTemplate
	MessageType int
	FieldValues map[int]FieldValue
}

func (m *BitmapMessage) Init() {
	m.FieldValues = make(map[int]FieldValue)

}

type Unpacker interface {
	Unpack(offset int, data []byte) (newOffset int, fieldData []byte, err error)
}
type Packer interface {
	Pack(data []byte) (msg []byte, err error)
}

type PackerUnpacker interface {
	Unpacker
	Packer
}

type FieldUnpacker interface {
	UnpackField(offset int, data []byte) (newOffset int, fieldValue FieldValue, err error)
}
type FieldPacker interface {
	PackField(fieldValue FieldValue) (msg []byte, err error)
}

type FieldPackerUnpacker interface {
	FieldUnpacker
	FieldPacker
}

type fieldType int

const (
	AlphaNumeric fieldType = iota
	Numeric
	AlphaNumericSpecial
	Alpha
	Binary
)

var fieldTypeLookup = map[fieldType]string{
	AlphaNumeric:        "an",
	Numeric:             "n",
	AlphaNumericSpecial: "ans",
	Alpha:               "a",
	Binary:              "b",
}

type Field interface {
	GetFieldNumber() int
	GetName() string
	GetType() fieldType
	GetLength() variableFieldLength
	GetSize() int
	FieldPackerUnpacker
}

type SubField interface {
	GetSubField(field int) (value string, err error)
}

type BitmapMessageField struct {
	FieldNumber int
	Name        string
	Type        fieldType
	Length      variableFieldLength //Type of field length
	Size        int                 //Maximum size of field
	PackerUnpacker
	FieldPackerUnpacker
}

type FieldValue struct {
	//	Field
	Value       string
	FieldValues map[int]FieldValue
	fmt.Stringer
}

func (f *FieldValue) String() string {
	return f.Value
}

func (f *BitmapMessageField) GetFieldNumber() int {
	return f.FieldNumber
}
func (f *BitmapMessageField) GetName() string {
	return f.Name
}
func (f *BitmapMessageField) GetType() fieldType {
	return f.Type
}
func (f *BitmapMessageField) GetLength() variableFieldLength {
	return f.Length
}
func (f *BitmapMessageField) GetSize() int {
	return f.Size
}

func (f *BitmapMessageField) ValidateValue(value string) (valid bool, err error) {
	return true, nil
}

func (f *BitmapMessageField) UnpackField(offset int, data []byte) (newOffset int, value FieldValue, err error) {
	newOffset, fieldData, err := f.Unpack(offset, data)
	value = FieldValue{Value: string(fieldData)}
	return newOffset, value, nil
	//return offset, "", errors.New("Undefined unpacking")
}
func (f *BitmapMessageField) PackField(value FieldValue) (data []byte, err error) {

	var fieldData []byte
	if f.GetLength() == Fixed {
		fieldData, err = f.Pack([]byte(util.LeftPad2Len(value.Value, " ", f.GetSize())))
	} else {
		fieldData, err = f.Pack([]byte(value.Value))
	}
	return fieldData, err
}

func (f *BitmapMessageTemplate) GetFieldDef(fieldNumber int) (field Field, err error) {
	field, exists := f.Fields[fieldNumber]
	if !exists {
		return nil, errors.New("Field not found")
	}
	return field, nil
}

func (m *BitmapMessage) SetString(fieldNr int, value string) {
	m.SetField(fieldNr, FieldValue{Value: value})
}
func (m *BitmapMessage) SetField(fieldNr int, value FieldValue) {
	if fieldNr > 64 && !m.IsFieldSet(1) {
		m.SetField(1, *new(FieldValue))
	}

	m.FieldValues[fieldNr] = value
	//TODO: Validate field?
}

func (m *BitmapMessage) SetSubField(fieldNr int, subFieldNr int, value string) {

	if _, ok := m.FieldValues[fieldNr]; !ok {
		m.SetField(fieldNr, FieldValue{FieldValues: make(map[int]FieldValue)})
	}

	m.FieldValues[fieldNr].FieldValues[subFieldNr] = FieldValue{Value: value}
}

func (m *BitmapMessage) GetSubField(fieldNr, subFieldNr int) (value string, isSet bool) {
	if _, ok := m.FieldValues[fieldNr]; !ok {
		return value, ok //Not found
	}
	if m.FieldValues[fieldNr].FieldValues == nil {
		//not a set subfield.
		return value, false
	}
	if fvalue, ok := m.FieldValues[fieldNr].FieldValues[subFieldNr]; !ok {
		return value, ok
	} else {
		return fvalue.Value, ok
	}

}
func (m *BitmapMessage) IsFieldSet(fieldNr int) bool {
	_, set := m.FieldValues[fieldNr]
	return set
}
func (m *BitmapMessage) GetField(fieldNr int) (value string, isSet bool) {
	fieldValue, set := m.FieldValues[fieldNr]
	return fieldValue.Value, set
}
func (m *BitmapMessage) GetString(fieldNr int) (value string, isSet bool) {
	return m.GetField(fieldNr)
}

func (m *BitmapMessage) CopyField(fieldNr int, msg Message) {
	v, set := msg.GetField(fieldNr)
	if set {
		m.SetString(fieldNr, v)
	}

}

func (m *BitmapMessage) GetMsgType() int {
	return m.MessageType
}

func (m *BitmapMessage) GetMsgTypeString() string {
	return util.LeftPad2Len(strconv.FormatInt(int64(m.GetMsgType()), 16), "0", 4)
}

func (m *BitmapMessage) Pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteString(m.GetMsgTypeString())

	//Generate the bitmap based on set fields.
	bytes, err := PackBitmapFields(m.FieldValues, m.BitmapMessageTemplate)
	if err != nil {
		return nil, err
	}
	buf.Write(bytes)
	return buf.Bytes(), nil
}

func PackBitmapFields(fieldValues map[int]FieldValue, tmpl *BitmapMessageTemplate) (data []byte, err error) {
	buf := new(bytes.Buffer)
	bufFields := new(bytes.Buffer)

	var setFields []int
	for m := range fieldValues {
		setFields = append(setFields, m)
	}
	sort.Ints(setFields)
	var bitmapSize int = 1
	if _, ok := fieldValues[1]; ok {
		bitmapSize = 2
	}
	bitMap := make([]uint64, bitmapSize)

	for _, i := range setFields {
		ui := uint64(i)
		//bitMap[i/64] |= 1 << (64 - (ui - (ui/64)*64))
		//bitMap[i/64] |= 1 << (64 - ui%64)
		bitMap[i/64] |= 1 << (ui%64 - 1) //We set each bit accordingly, then reverse the whole map when packing

		if i > 1 {
			f, ok := tmpl.Fields[i]
			if !ok {
				fmt.Println("Template missing: ", i)
			}
			fieldBytes, err := f.PackField(fieldValues[i])
			if err != nil {
				fmt.Println("error packing field ", i, " ", err)
				continue
			}
			binary.Write(bufFields, binary.LittleEndian, fieldBytes)
		}

	}

	for _, b := range bitMap {
		binary.Write(buf, binary.BigEndian, util.ReverseUint64Bits(b))
	}
	//Write fields
	bufFields.WriteTo(buf)
	return buf.Bytes(), nil
}

func (m *BitmapMessage) SetMsgType(msgType int) {
	m.MessageType = msgType
}

//func (f *BitmapMessageField) Unpacker(offset int, data []byte) (newOffset int, fieldValue string, err error) {
//	return 0, "", errors.New("Error")
//}

func CreateFields(fields ...Field) map[int]Field {
	fieldMap := make(map[int]Field)
	for _, f := range fields {
		fieldMap[f.GetFieldNumber()] = f
	}
	return fieldMap
}

func BitmapUnpack(data []byte, tmpl MessageTemplate, msg Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	if msg == nil {
		return errors.New("Message must not be nil")
	}
	if len(data) < 12 {
		return errors.New("Invalid message size. Must be at least 12 bytes")
	}
	//First 4 bytes are the message type header.
	msgTypeStr := string(data[0:4])
	msgType, err := strconv.ParseInt(msgTypeStr, 16, 0)

	if err != nil {
		return errors.New(fmt.Sprint("Message type not numeric; ", err))
	}
	msg.SetMsgType(int(msgType))
	i := 4
	bitmap := data[i : i+8]
	i = i + 8
	if bitmap[0]&0x80 > 0 { //Extended bitmap.
		bitmap = append(bitmap, data[i:i+8]...)
		i = i + 8
	}

	boolBitmap := util.GetBitmap(bitmap)

	for fieldNr, b := range boolBitmap {
		fieldNr = fieldNr + 1
		if fieldNr == 1 {
			continue
		}
		if b {
			field, err := tmpl.GetFieldDef(fieldNr)
			if err != nil {
				fmt.Println("Error unpacking field ", fieldNr, " : ", err)
				return err
			}
			var fieldValue FieldValue
			i, fieldValue, err = field.UnpackField(i, data)

			if err != nil {
				fmt.Println("Error unpacking field ", fieldNr, " : ", err)
				return err
			}
			msg.SetField(fieldNr, fieldValue)
		}
	}
	return nil
}

func (m *BitmapMessage) String() string {
	buf := bytes.NewBufferString("")

	buf.WriteString(util.LeftPad2Len(strconv.FormatInt(int64(m.GetMsgType()), 16), "0", 4))
	buf.WriteString(":\n")
	/*
	   [LLVAR  n    ..19 016] 002 [4300000000008267]
	   [Fixed  n       6 006] 003 [350000]
	   [None   n         012] 004 [000000004000]
	   [Fixed  n      10 010] 007 [0602195310]
	   [Fixed  n       6 006] 011 [001594]
	*/
	formatFieldsToString(buf, m.BitmapMessageTemplate, m.FieldValues, "")
	return buf.String()
}

func formatFieldsToString(buf *bytes.Buffer, tmpl MessageTemplate, values map[int]FieldValue, fieldPrefix string) {
	var setFields []int
	for m := range values {
		setFields = append(setFields, m)
	}
	sort.Ints(setFields)

	for _, i := range setFields {
		if i > 1 {
			f, err := tmpl.GetFieldDef(i)
			if err != nil {
				fmt.Println("Template missing: ", i)
			}
			fieldValue := values[i]
			if fieldValue.FieldValues != nil {
				mt, ok := f.(MessageTemplate)
				if ok {
					formatFieldsToString(buf, mt, fieldValue.FieldValues, fmt.Sprint(strconv.Itoa(i), "."))
				}
			} else {
				formatFieldToString(buf, f, fieldValue, fieldPrefix)
			}

		}
	}
}

func formatFieldToString(buf *bytes.Buffer, field Field, fieldValue FieldValue, fieldNrPrefix string) {

	buf.WriteString("\t[")
	buf.WriteString(util.RightPad2Len(field.GetLength().String(), " ", 8))
	buf.WriteString(util.RightPad2Len(fieldTypeLookup[field.GetType()], " ", 5))
	buf.WriteString(util.LeftPad2Len(strconv.Itoa(field.GetSize()), " ", 4))
	buf.WriteString(" ")

	//fieldValue := m.FieldValues[i]

	buf.WriteString(util.LeftPad2Len(strconv.Itoa(len(fieldValue.String())), "0", 3))
	buf.WriteString("] ")
	buf.WriteString(util.RightPad2Len(fmt.Sprint(fieldNrPrefix, util.LeftPad2Len(strconv.Itoa(field.GetFieldNumber()), "0", 3)), " ", 7))

	buf.WriteString(" [")
	buf.WriteString(fieldValue.String())
	buf.WriteString("]")
	buf.WriteString("\n")

}
