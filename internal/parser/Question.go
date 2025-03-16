package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type QClass uint16

const (
	IN  QClass = 1
	CS  QClass = 2
	CH  QClass = 3
	HS  QClass = 4
	ANY QClass = 255
)

type QType uint16

const (
	A  QType = 1
	NS QType = 2
	MD QType = 3
	MF QType = 4
)

type Question struct {
	Labels []string
	Type   QType
	Class  QClass
}

func ParseQuestion(buffer *bytes.Buffer) Question {
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
	buf := make([]byte, 2)
	buffer.Read(buf)
	qtype := binary.BigEndian.Uint16(buf)

	buf = make([]byte, 2)
	buffer.Read(buf)
	qclass := binary.BigEndian.Uint16(buf)

	return Question{
		Labels: labels,
		Type:   QType(qtype),
		Class:  QClass(qclass),
	}
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
