package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

type QType uint16
type QClass uint16
type OPCODE uint8

const (
	QUERY OPCODE = iota
	IQUERY
	STATUS
)

type RCODE uint8

const (
	NO_ERROR RCODE = iota
	FORMAT_ERROR
	SERVER_FAILURE
	NAME_ERROR
	NOT_IMPLEMENTED
	REFUSED
)

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

func ParseHeader(buf []byte) Header {
	buffer := bytes.NewBuffer(buf)
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

type Question struct {
	Labels []string
	Type   QType
	Class  QClass
}

func (question Question) String() string {
	return fmt.Sprintf("question: %s type: %d class: %d", strings.Join(question.Labels, "."), question.Type, question.Class)
}

func (question Question) ToBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	for i := 0; i < len(question.Labels); i++ {
		size := byte(len(question.Labels[i]))
		if err := binary.Write(buf, binary.BigEndian, size); err != nil {
			return nil, err
		}
		buf.Write([]byte(question.Labels[i]))
	}
	buf.WriteByte(0)
	if err := binary.Write(buf, binary.BigEndian, question.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, question.Class); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Answer struct {
	Labels []string
	Type   QType
	Class  QClass
	TTL    uint32
	Data   []byte
}

func (answer Answer) ToBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	for i := 0; i < len(answer.Labels); i++ {
		size := byte(len(answer.Labels[i]))
		if err := binary.Write(buf, binary.BigEndian, size); err != nil {
			return nil, err
		}
		buf.Write([]byte(answer.Labels[i]))
	}
	buf.WriteByte(0)
	if err := binary.Write(buf, binary.BigEndian, answer.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, answer.Class); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, answer.TTL); err != nil {
		return nil, err
	}
	length := uint16(len(answer.Data))
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return nil, err
	}
	buf.Write(answer.Data)
	return buf.Bytes(), nil
}

func main() {
	conn, err := net.ListenPacket("udp", ":53")
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()
	for {
		request := make([]byte, 1024)
		_, addr, _ := conn.ReadFrom(request)
		requestheader, resquestquestion := ParseMessage(request)

		response := new(bytes.Buffer)
		responseHeader := Header{
			ID:            requestheader.ID,
			QuestionCount: 1,
			AnswerCount:   1,
		}
		request, _ = responseHeader.ToBinary()
		response.Write(request)

		question, _ := resquestquestion.ToBinary()
		response.Write(question)

		answer := Answer{
			Labels: resquestquestion.Labels,
			Type:   resquestquestion.Type,
			Class:  resquestquestion.Class,
			TTL:    3600,
			Data:   []byte{1, 1, 1, 1},
		}
		answerBuf, _ := answer.ToBinary()
		response.Write(answerBuf)
		conn.WriteTo(response.Bytes(), addr)
	}

}

func ParseMessage(buf []byte) (Header, Question) {
	header := ParseHeader(buf[:4*3])

	question := buf[4*3:]

	buffer := bytes.NewBuffer(question)
	labels := []string{}
	for {
		length, err := buffer.ReadByte()
		if length == 0 {
			break
		}
		if err != nil {
			os.Exit(1)
		}
		label := make([]byte, length)
		_, err = buffer.Read(label)
		if err != nil {
			os.Exit(1)
		}
		labels = append(labels, string(label))
	}
	buf = make([]byte, 2)
	buffer.Read(buf)
	qtype := binary.BigEndian.Uint16(buf)

	buf = make([]byte, 2)
	buffer.Read(buf)

	return header, Question{
		Labels: labels,
		Type:   QType(qtype),
		Class:  QClass(qtype),
	}
}
