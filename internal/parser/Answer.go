package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Answer struct {
	Labels []string
	Type   QType
	Class  QClass
	TTL    uint32
	Data   []byte
}

func ParseAnswer(buffer *MessageBuffer) Answer {
	labels, _ := buffer.ReadLabels()
	t := QType(buffer.ReadUint16())
	class := buffer.ReadUint16()
	ttl := buffer.ReadUint32()
	length := buffer.ReadUint16()
	var data []byte

	switch t {
	case NS:
		l, _ := buffer.ReadLabels()
		buf := new(bytes.Buffer)
		for i := 0; i < len(l); i++ {
			size := byte(len(l[i]))
			binary.Write(buf, binary.BigEndian, size)
			buf.Write([]byte(l[i]))
		}
		buf.WriteByte(0)
		data = buf.Bytes()

	default:
		data = make([]byte, length)
		buffer.Read(data)
	}
	return Answer{
		Labels: labels,
		Type:   t,
		Class:  QClass(class),
		TTL:    ttl,
		Data:   data,
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
