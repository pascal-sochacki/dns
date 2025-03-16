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

func ParseHeader(buffer *bytes.Buffer) Header {
	id := make([]byte, 2)
	buffer.Read(id)
	parsedId := binary.BigEndian.Uint16(id)

	flags := make([]byte, 2)
	buffer.Read(flags)
	flagsNumber := binary.BigEndian.Uint16(flags)

	opcode := QUERY

	if (flagsNumber & (1 << 11)) != 0 {
		opcode = IQUERY
	} else if (flagsNumber & (1 << 12)) != 0 {
		opcode = STATUS
	}

	rcode := flagsNumber & (0b00000000_00001111)

	count := make([]byte, 2)
	buffer.Read(count)
	questionCount := binary.BigEndian.Uint16(count)

	buffer.Read(count)
	answerCount := binary.BigEndian.Uint16(count)

	buffer.Read(count)
	nsCount := binary.BigEndian.Uint16(count)

	buffer.Read(count)
	arCount := binary.BigEndian.Uint16(count)

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
