package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type Answer struct {
	Labels []string
	Type   QType
	Class  QClass
	TTL    uint32
	Data   []byte
}

func ParseAnswer(buffer *bytes.Buffer) Answer {
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
		fmt.Printf("%08b ", length)
		for i := 0; i < len(label); i++ {
			fmt.Printf("%08b ", label[i])

		}
		labels = append(labels, string(label))
	}

	buf := make([]byte, 2)

	buffer.Read(buf)
	qtype := binary.BigEndian.Uint16(buf)

	buffer.Read(buf)
	class := binary.BigEndian.Uint16(buf)

	buf = make([]byte, 4)
	buffer.Read(buf)
	ttl := binary.BigEndian.Uint32(buf)

	buf = make([]byte, 2)
	buffer.Read(buf)
	length := binary.BigEndian.Uint16(buf)

	buf = make([]byte, length)
	buffer.Read(buf)

	return Answer{
		Labels: labels,
		TTL:    ttl,
		Type:   QType(qtype),
		Class:  QClass(class),
		Data:   buf,
	}

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

func (answer Answer) String() string {
	full := strings.Join(answer.Labels, ".")
	return fmt.Sprintf("%s ttl: %d type: %d class: %d data: %s", full, answer.TTL, answer.Type, answer.Class, string(answer.Data))

}
