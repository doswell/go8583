package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/doswell/go8583"
	"github.com/doswell/go8583/util"
)

func main() {

	//Create an empty new field
	message := NewIso8583Message()
	message.SetString(2, "2342434232")
	message.SetString(3, "002000")

	pan, ok := message.GetString(2)
	if ok {
		fmt.Println("Recorded pan is", pan)
	}

}

//Iso8583 Example standard Is8583 message layout per
type Iso8583 struct {
	*go8583.BitmapMessage
}

var iso8583MsgTemplate = (&go8583.BitmapMessageTemplate{
	Header: []go8583.Field{},
	Fields: go8583.CreateFields(
		go8583.NewFixedField(1, "extendedBitMap", 8, go8583.Binary),
		go8583.NewLlVarField(2, "pan", 2, go8583.AlphaNumeric),
		go8583.NewFixedField(3, "processingCode", 6, go8583.Numeric),
		go8583.NewFixedField(4, "amountTransaction", 12, go8583.Numeric),
		go8583.NewFixedField(5, "amountSettlement", 12, go8583.Numeric),
		go8583.NewFixedField(6, "amountcardholderBilling", 12, go8583.Numeric),
		go8583.NewFixedField(7, "transmissionDateTime", 10, go8583.Numeric),
		go8583.NewFixedField(9, "conversionRateSettlement", 8, go8583.Numeric),
		go8583.NewFixedField(11, "traceNumber", 6, go8583.Numeric),
		go8583.NewFixedField(12, "localTranTime", 6, go8583.Numeric),
		go8583.NewFixedField(13, "localTranDate", 4, go8583.Numeric),
		go8583.NewFixedField(14, "expirationDate", 4, go8583.Numeric),
		go8583.NewFixedField(15, "settlementDate", 4, go8583.Numeric),
		go8583.NewFixedField(18, "merchantType", 4, go8583.Numeric),
		go8583.NewFixedField(22, "posEntryMode", 3, go8583.Numeric),
		go8583.NewFixedField(23, "cardSequenceNumber", 3, go8583.Numeric),
		go8583.NewFixedField(25, "posConditionCode", 2, go8583.Numeric),
		go8583.NewFixedField(26, "posPinCaptureCode", 2, go8583.Numeric),
		go8583.NewFixedField(28, "tranFee", 9, go8583.AlphaNumeric),
		go8583.NewFixedField(30, "settleFee", 9, go8583.AlphaNumeric),
		go8583.NewLlVarField(32, "acquiingInstId", 11, go8583.Numeric),
		go8583.NewLlVarField(35, "track2", 37, go8583.AlphaNumericSpecial),
		go8583.NewFixedField(37, "retrievalReferneceNumber", 12, go8583.AlphaNumericSpecial),
		go8583.NewFixedField(38, "authoizationCode", 6, go8583.AlphaNumericSpecial),
		go8583.NewFixedField(39, "responseCode", 2, go8583.AlphaNumeric),
		go8583.NewFixedField(40, "serviceRestrictionCode", 3, go8583.AlphaNumeric),

		go8583.NewFixedField(41, "terminalId", 8, go8583.AlphaNumericSpecial),
		go8583.NewFixedField(42, "cardAcceptorId", 15, go8583.AlphaNumericSpecial),
		go8583.NewFixedField(43, "cardAcceptorNameLoc", 40, go8583.AlphaNumericSpecial),
		//   [LLVAR  ans  ..25 003] 044 [018]
		go8583.NewLlVarField(44, "additionalRspData", 25, go8583.AlphaNumericSpecial),
		go8583.NewLllVarField(48, "additionalData", 999, go8583.AlphaNumericSpecial),
		go8583.NewFixedField(49, "currencyCodeTran", 3, go8583.Numeric),
		go8583.NewLllVarField(54, "extendedAmounts", 120, go8583.AlphaNumeric),
		go8583.NewLllVarField(55, "iccData", 300, go8583.Binary),
		go8583.NewLllVarField(57, "authorizationLifecycleCode", 3, go8583.Numeric),
		go8583.NewLllVarField(59, "echoData", 500, go8583.AlphaNumericSpecial),
		go8583.NewFixedField(70, "networkMgmtCode", 3, go8583.Numeric),
		go8583.NewFixedField(90, "originalDataElements", 42, go8583.AlphaNumeric),
		go8583.NewFixedField(91, "fileUpdateCode", 1, go8583.AlphaNumeric),
		// [Fixed  an*    42 042] 095 [000000010000000000010000C00000000C00000000]
		go8583.NewFixedField(95, "replacementAmounts", 42, go8583.AlphaNumeric),
		go8583.NewLlVarField(100, "receivingInstId", 11, go8583.Numeric),
		go8583.NewLlVarField(101, "fileName", 17, go8583.AlphaNumeric),
		go8583.NewLllVarField(123, "customField", 999, go8583.AlphaNumeric),

		//Custom field with field presence denoted by a bitmap field.
		go8583.NewBitmapField(127, "customBitmapSubFields",
			go8583.NewVariableFieldPackerUnpacker(int(go8583.LlllllVar)),
			[]go8583.Field{
				go8583.NewLlVarField(2, "customField1", 60, go8583.AlphaNumericSpecial),
				go8583.NewLlVarField(3, "customField2", 60, go8583.AlphaNumericSpecial),
				go8583.NewLlVarField(4, "customField3", 60, go8583.AlphaNumericSpecial),
				go8583.NewLlVarField(5, "customField4", 60, go8583.AlphaNumericSpecial),
				go8583.NewLlVarField(10, "customField5", 60, go8583.AlphaNumericSpecial),
				go8583.NewLlVarField(20, "customField6", 60, go8583.AlphaNumericSpecial),
			}),
	),
})

func NewIso8583Message() *Iso8583 {
	msg := new(Iso8583)
	msg.BitmapMessage = new(go8583.BitmapMessage)
	msg.BitmapMessage.BitmapMessageTemplate = iso8583MsgTemplate
	msg.BitmapMessage.Init()
	return msg
}

func NewIso8583MessageFromByte(data []byte) (*Iso8583, error) {
	msg := NewIso8583Message()

	err := go8583.BitmapUnpack(data, iso8583MsgTemplate, msg)
	return msg, err
}

func NewIso8583_0800Message(ntwkCode string, trace int) *Iso8583 {
	msg := NewIso8583Message()
	msg.SetMsgType(0x0800)
	msg.SetString(7, time.Now().UTC().Format("0102150405"))
	msg.SetString(11, util.LeftPad2Len(strconv.Itoa(trace), "0", 6))
	msg.SetString(12, time.Now().Format("150405"))
	msg.SetString(13, time.Now().Format("0102"))
	msg.SetString(70, ntwkCode)
	return msg
}
