package go8583

type variableFieldLength int

//go:generate stringer -type=variableFieldLength
const (
	Fixed variableFieldLength = iota
	LVar
	LlVar
	LllVar
	LlllVar
	LllllVar
	LlllllVar
)
