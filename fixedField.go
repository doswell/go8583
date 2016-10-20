package go8583

import "errors"

type fixedField struct {
	Size int
}

func (f *fixedField) Unpack(offset int, data []byte) (newOffset int, fieldData []byte, err error) {
	fieldEnd := offset + f.Size
	if fieldEnd > len(data) {
		return 0, fieldData, errors.New("Attempt to read passed end of data")
	}
	fieldData = data[offset:fieldEnd]
	return fieldEnd, fieldData, nil
}

func (f *fixedField) Pack(fieldData []byte) (data []byte, err error) {
	if len(fieldData) != f.Size {
		return nil, errors.New("Field data size not equal to total field size for a fixed field.")
	}
	//	if isOk, err := f.ValidateValue(value.Value); !isOk {
	//		return nil, err
	//	}
	//	fmt.Println("fixed field packing : ", []byte(value.Value))
	return fieldData, nil
}

func NewFixedField(bitNumber int, name string, size int, fieldType fieldType) Field {
	return &BitmapMessageField{bitNumber, name, fieldType, Fixed, size, NewFixedFieldPackerUnpacker(size), nil}
}

func NewFixedFieldPackerUnpacker(size int) PackerUnpacker {
	return &fixedField{Size: size}
}
