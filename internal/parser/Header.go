package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type OPCODE uint8

const (
	QUERY OPCODE = iota
	IQUERY
	STATUS
)

const (
	NO_ERROR RCODE = iota
	FORMAT_ERROR
	SERVER_FAILURE
	NAME_ERROR
	NOT_IMPLEMENTED
	REFUSED
)

type RCODE uint8

type Header struct {
	ID                  uint16
	IsQuery             bool
	OPCODE              OPCODE
	AuthoritativeAnswer bool
	TrunCation          bool
	RecursionDesired    bool
	RecursionAvailable  bool
	ResponseCode        RCODE

	QuestionCount uint16
	AnswerCount   uint16
	NSCount       uint16
	ARCount       uint16
}

func ParseHeader(buffer *MessageBuffer) Header {
	parsedId := buffer.ReadUint16()
	flagsNumber := buffer.ReadUint16()

	opcode := QUERY

	if (flagsNumber & (1 << 11)) != 0 {
		opcode = IQUERY
	} else if (flagsNumber & (1 << 12)) != 0 {
		opcode = STATUS
	}

	rcode := flagsNumber & (0b00000000_00001111)

	questionCount := buffer.ReadUint16()
	answerCount := buffer.ReadUint16()
	nsCount := buffer.ReadUint16()
	arCount := buffer.ReadUint16()

	return Header{
		ID:            parsedId,
		OPCODE:        opcode,
		ResponseCode:  RCODE(rcode),
		IsQuery:       (flagsNumber & (1 << 15)) == 0,
		QuestionCount: questionCount,
		AnswerCount:   answerCount,
		ARCount:       arCount,
		NSCount:       nsCount,
	}
}

func (header Header) String() string {
	return fmt.Sprintf("id: %d, qcount: %d acount: %d, nscount: %d, arcount: %d", header.ID, header.QuestionCount, header.AnswerCount, header.NSCount, header.ARCount)
}

func (header Header) ToBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, header.ID); err != nil {
		return nil, err
	}

	var flags uint16
	if !header.IsQuery {
		flags |= uint16(1 << 15)
	}
	flags |= uint16(header.OPCODE) << 11
	flags |= uint16(header.ResponseCode)

	if err := binary.Write(buf, binary.BigEndian, flags); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, header.QuestionCount); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, header.AnswerCount); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, header.NSCount); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, header.ARCount); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
