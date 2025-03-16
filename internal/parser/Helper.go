package parser

import "encoding/binary"

type MessageBuffer struct {
	buf []byte
	off int
}

func NewLookBackBuffer(b []byte) *MessageBuffer {
	return &MessageBuffer{
		buf: []byte(b),
	}
}

func (r *MessageBuffer) Read(b []byte) (n int, err error) {
	n = copy(b, r.buf[r.off:])
	r.off += n
	return n, nil
}

func (r *MessageBuffer) ReadUint16() (n uint16) {
	count := make([]byte, 2)
	r.Read(count)
	return binary.BigEndian.Uint16(count)
}
func (r *MessageBuffer) ReadUint32() (n uint32) {
	count := make([]byte, 4)
	r.Read(count)
	return binary.BigEndian.Uint32(count)
}

func (r *MessageBuffer) ReadByte() (n byte, err error) {
	result := r.buf[r.off]
	r.off += 1
	return result, nil
}

func (r *MessageBuffer) ReadLabels() ([]string, error) {
	labels := []string{}
	for {
		length, err := r.ReadByte()

		if (length & 0b11000000) == 0b11000000 {
			octets, _ := r.ReadByte()
			offset := binary.BigEndian.Uint16([]byte{
				(length & 0b00111111),
				octets,
			})
			current := r.off
			r.off = int(offset)
			pointerLabels, _ := r.ReadLabels()
			r.off = current

			labels = append(labels, pointerLabels...)
			return labels, nil
		}

		if length == 0 {
			break
		}
		if err != nil {
			return labels, err
		}
		label := make([]byte, length)
		_, err = r.Read(label)
		if err != nil {
			return labels, err
		}
		labels = append(labels, string(label))
	}
	return labels, nil
}
