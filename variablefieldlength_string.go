// Code generated by "stringer -type=variableFieldLength"; DO NOT EDIT

package go8583

import "fmt"

const _variableFieldLength_name = "FixedLVarLlVarLllVarLlllVarLllllVarLlllllVar"

var _variableFieldLength_index = [...]uint8{0, 5, 9, 14, 20, 27, 35, 44}

func (i variableFieldLength) String() string {
	if i < 0 || i >= variableFieldLength(len(_variableFieldLength_index)-1) {
		return fmt.Sprintf("variableFieldLength(%d)", i)
	}
	return _variableFieldLength_name[_variableFieldLength_index[i]:_variableFieldLength_index[i+1]]
}
