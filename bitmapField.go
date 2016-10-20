package go8583

import (
	"doswell/util"
	"errors"
	"fmt"
)

type bitmapField struct {
	*BitmapMessageField
	*BitmapMessageTemplate
}

func NewBitmapField(fieldNumber int, name string, fieldPackerUnpacker PackerUnpacker, fields []Field) Field {
	field := &bitmapField{&BitmapMessageField{fieldNumber, name, Binary, LlVar, 0, fieldPackerUnpacker, nil},
		new(BitmapMessageTemplate),
	}
	field.Fields = CreateFields(fields...)

	return field
}

func (f *bitmapField) UnpackField(offset int, data []byte) (newOffset int, value FieldValue, err error) {
	newOffset, fieldData, err := f.Unpack(offset, data)
	//Get the bitmap.
	bitmapData := fieldData[0:8]
	bitmap := util.GetBitmap(bitmapData)

	bitmapFieldValue := new(FieldValue)
	bitmapFieldValue.FieldValues = make(map[int]FieldValue)
	fieldOffset := 8

	for i, b := range bitmap {
		i++
		if b {
			field, err := f.GetFieldDef(i)
			if err != nil {
				return newOffset, value, errors.New(fmt.Sprint("Undefined field ", i, " in template"))
			}
			var fValue FieldValue

			fieldOffset, fValue, err = field.UnpackField(fieldOffset, fieldData)
			if err != nil {
				return newOffset, value, err
			}

			bitmapFieldValue.FieldValues[i] = fValue
		}
	}

	return newOffset, *bitmapFieldValue, nil
	//return offset, "", errors.New("Undefined unpacking")
}
func (f *bitmapField) PackField(value FieldValue) (data []byte, err error) {

	bitmapBytes, err := PackBitmapFields(value.FieldValues, f.BitmapMessageTemplate)
	if err != nil {
		return data, err
	}

	return f.Pack(bitmapBytes)
}
