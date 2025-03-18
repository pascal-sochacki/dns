package main

import (
	"testing"

	"github.com/pascal-sochacki/dns/internal/parser"
)

func TestParseMessage(t *testing.T) {
	tests := []struct {
		input              []byte
		header             parser.Header
		question           parser.Question
		nameserverResouece []parser.Answer
	}{
		{
			input: []byte{
				0b00000000, 0b00000001,
				0b00000000, 0b00000000,
				0b00000000, 0b00000001,
				0b00000000, 0b00000000,
				0b00000000, 0b00000000,
				0b00000000, 0b00000000,

				4, 'b', 'l', 'o', 'g',
				7, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
				3, 'c', 'o', 'm',
				0,

				0b00000000,
				0b00000001,

				0b00000000,
				0b00000001,
			},
			header: parser.Header{
				ID:            1,
				IsQuery:       true,
				ResponseCode:  parser.NO_ERROR,
				OPCODE:        parser.QUERY,
				QuestionCount: 1,
				AnswerCount:   0,
				NSCount:       0,
				ARCount:       0,
			},
			question: parser.Question{
				Class: parser.IN,
				Type:  parser.A,
				Labels: []string{
					"blog",
					"example",
					"com",
				},
			},
		},
		{
			input: []byte{
				0b00000000, 0b00001100,
				0b10000000, 0b00000000,
				0b00000000, 0b00000001,
				0b00000000, 0b00000000,
				0b00000000, 0b00000101,
				0b00000000, 0b00001001,

				0b00000010, 0b01100101, 0b01110101,
				0b00000000,
				0b00000000, 0b00000001,
				0b00000000, 0b00000001,

				// compression with offset
				0b11000000, 0b00001100,

				0b00000000, 0b00000010, 0b00000000, 0b00000001,
				0b00000000, 0b00000010, 0b10100011, 0b00000000,
				0b00000000, 0b00001000, 0b00000001, 0b01110111,
				0b00000011, 0b01100100, 0b01101110, 0b01110011,
				0b11000000, 0b00001100, 0b11000000, 0b00001100,
				0b00000000, 0b00000010, 0b00000000, 0b00000001,
				0b00000000, 0b00000010, 0b10100011, 0b00000000,
				0b00000000, 0b00000100, 0b00000001, 0b01111000,
				0b11000000, 0b00100010, 0b11000000, 0b00001100,
				0b00000000, 0b00000010, 0b00000000, 0b00000001,
				0b00000000, 0b00000010, 0b10100011, 0b00000000,
				0b00000000, 0b00000100, 0b00000001, 0b01111001,
				0b11000000, 0b00100010, 0b11000000, 0b00001100,
				0b00000000, 0b00000010, 0b00000000, 0b00000001,
				0b00000000, 0b00000010, 0b10100011, 0b00000000,
				0b00000000, 0b00000101, 0b00000010, 0b01100010,
				0b01100101, 0b11000000, 0b00100010, 0b11000000,
				0b00001100, 0b00000000, 0b00000010, 0b00000000,
				0b00000001, 0b00000000, 0b00000010, 0b10100011,
				0b00000000, 0b00000000, 0b00000101, 0b00000010,
				0b01110011, 0b01101001, 0b11000000, 0b00100010,
				0b11000000, 0b00100000, 0b00000000, 0b00000001,
				0b00000000, 0b00000001, 0b00000000, 0b00000010,
				0b10100011, 0b00000000, 0b00000000, 0b00000100,
				0b11000010, 0b00000000, 0b00011001, 0b00011100,
				0b11000000, 0b00110100, 0b00000000, 0b00000001,
				0b00000000, 0b00000001, 0b00000000, 0b00000010,
				0b10100011, 0b00000000, 0b00000000, 0b00000100,
				0b10111001, 0b10010111, 0b10001101, 0b00000001,
				0b11000000, 0b01000100, 0b00000000, 0b00000001,
				0b00000000, 0b00000001, 0b00000000, 0b00000010,
				0b10100011, 0b00000000, 0b00000000, 0b00000100,
				0b11000010, 0b10010010, 0b01101010, 0b01011010,
				0b11000000, 0b01010100, 0b00000000, 0b00000001,
				0b00000000, 0b00000001, 0b00000000, 0b00000010,
				0b10100011, 0b00000000, 0b00000000, 0b00000100,
				0b10010101, 0b00100110, 0b00000001, 0b00011010,
				0b11000000, 0b01100101, 0b00000000, 0b00000001,
				0b00000000, 0b00000001, 0b00000000, 0b00000010,
				0b10100011, 0b00000000, 0b00000000, 0b00000100,
				0b11000001, 0b00000010, 0b11011101, 0b00111110,
				0b11000000, 0b00100000, 0b00000000, 0b00011100,
				0b00000000, 0b00000001, 0b00000000, 0b00000010,
				0b10100011, 0b00000000, 0b00000000, 0b00010000,
				0b00100000, 0b00000001, 0b00000110, 0b01111000,
				0b00000000, 0b00100000, 0b00000000, 0b00000000,
				0b00000000, 0b00000000, 0b00000000, 0b00000000,
				0b00000000, 0b00000000, 0b00000000, 0b00101000,
				0b11000000, 0b00110100, 0b00000000, 0b00011100,
				0b00000000, 0b00000001, 0b00000000, 0b00000010,
				0b10100011, 0b00000000, 0b00000000, 0b00010000,
				0b00101010, 0b00000010, 0b00000101, 0b01101000,
				0b11111110, 0b00000000, 0b00000000, 0b00000000,
				0b00000000, 0b00000000, 0b00000000, 0b00000000,
				0b00000000, 0b00000000, 0b01100101, 0b01110101,
				0b11000000, 0b01000100, 0b00000000, 0b00011100,
				0b00000000, 0b00000001, 0b00000000, 0b00000010,
				0b10100011, 0b00000000, 0b00000000, 0b00010000,
				0b00100000, 0b00000001, 0b00000110, 0b01111100,
				0b00010000, 0b00010000, 0b00000000, 0b00100011,
				0b00000000, 0b00000000, 0b00000000, 0b00000000,
				0b00000000, 0b00000000, 0b00000000, 0b01010011,
				0b11000000, 0b01100101, 0b00000000, 0b00011100,
				0b00000000, 0b00000001, 0b00000000, 0b00000010,
				0b10100011, 0b00000000, 0b00000000, 0b00010000,
				0b00100000, 0b00000001, 0b00010100, 0b01110000,
				0b10000000, 0b00000000, 0b00000001, 0b00000000,
				0b00000000, 0b00000000, 0b00000000, 0b00000000,
				0b00000000, 0b00000000, 0b00000000, 0b01100010,
				0b00000000, 0b00000000, 0b00000000, 0b00000000,
				0b00000000, 0b00000000, 0b00000000, 0b00000000,
				0b00000000,
			},
			header: parser.Header{
				ID:            12,
				QuestionCount: 1,
				NSCount:       5,
				ARCount:       9,
			},
			question: parser.Question{
				Class:  parser.IN,
				Type:   parser.A,
				Labels: []string{"eu"},
			},
			nameserverResouece: []parser.Answer{
				{
					Labels: []string{"eu"},
					Type:   parser.NS,
					Class:  parser.IN,
					TTL:    172800,
					Data:   []byte{1, 'w', 3, 'd', 'n', 's', 2, 'e', 'u', 0},
				},
				{
					Labels: []string{"eu"},
					Type:   parser.NS,
					Class:  parser.IN,
					TTL:    172800,
					Data:   []byte{1, 'x', 3, 'd', 'n', 's', 2, 'e', 'u', 0},
				},
				{
					Labels: []string{"eu"},
					Type:   parser.NS,
					Class:  parser.IN,
					TTL:    172800,
					Data:   []byte{1, 'y', 3, 'd', 'n', 's', 2, 'e', 'u', 0},
				},
				{
					Labels: []string{"eu"},
					Type:   parser.NS,
					Class:  parser.IN,
					TTL:    172800,
					Data:   []byte{2, 'b', 'e', 3, 'd', 'n', 's', 2, 'e', 'u', 0},
				},
				{
					Labels: []string{"eu"},
					Type:   parser.NS,
					Class:  parser.IN,
					TTL:    172800,
					Data:   []byte{2, 's', 'i', 3, 'd', 'n', 's', 2, 'e', 'u', 0},
				},
			},
		},
	}
	for _, test := range tests {
		header, question, ns := ParseMessage(test.input)
		compareHeaders(t, header, test.header)
		compareQuestion(t, question, test.question)

		for i := 0; i < int(header.NSCount); i++ {
			compareAnswer(t, ns[i], test.nameserverResouece[i])
		}

	}
}

func compareHeaders(t *testing.T, is parser.Header, expect parser.Header) {
	t.Helper()
	if is.ID != expect.ID {
		t.Fatalf("id dont match is %d wanted %d", is.ID, expect.ID)
	}
	if is.QuestionCount != expect.QuestionCount {
		t.Fatalf("question count dont match is %d wanted %d", is.QuestionCount, expect.QuestionCount)
	}
	if is.AnswerCount != expect.AnswerCount {
		t.Fatalf("answer count dont match is %d wanted %d", is.AnswerCount, expect.AnswerCount)
	}
	if is.NSCount != expect.NSCount {
		t.Fatalf("ns count dont match is %d wanted %d", is.NSCount, expect.NSCount)
	}
	if is.ARCount != expect.ARCount {
		t.Fatalf("ar count dont match is %d wanted %d", is.ARCount, expect.ARCount)
	}
	if is.IsQuery != expect.IsQuery {
		t.Fatalf("is query dont match is %t wanted %t", is.IsQuery, expect.IsQuery)
	}
	if is.OPCODE != expect.OPCODE {
		t.Fatalf("op code dont match is %d wanted %d", is.OPCODE, expect.OPCODE)
	}
	if is.ResponseCode != expect.ResponseCode {
		t.Fatalf("response code dont match is %d wanted %d", is.ResponseCode, expect.ResponseCode)
	}
}

func TestHeaderToBinary(t *testing.T) {
	tests := []struct {
		expect []byte
		input  parser.Header
	}{
		{
			expect: []byte{
				0b00000000, 0b00000001,
				0b00000000, 0b00000000,
				0b00000000, 0b00000001,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
			},
			input: parser.Header{
				ID:            1,
				IsQuery:       true,
				ResponseCode:  parser.NO_ERROR,
				OPCODE:        parser.QUERY,
				QuestionCount: 1,
				AnswerCount:   2,
				NSCount:       3,
				ARCount:       4,
			},
		},
		{
			expect: []byte{
				0b00000000, 0b00000001,
				0b10001000, 0b00000001,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			input: parser.Header{
				ID:            1,
				IsQuery:       false,
				ResponseCode:  parser.FORMAT_ERROR,
				OPCODE:        parser.IQUERY,
				QuestionCount: 2,
				AnswerCount:   3,
				NSCount:       4,
				ARCount:       5,
			},
		},
		{
			expect: []byte{
				0b00000000, 0b00000001,
				0b10010000, 0b00000010,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			input: parser.Header{
				ID:            1,
				IsQuery:       false,
				ResponseCode:  parser.SERVER_FAILURE,
				OPCODE:        parser.STATUS,
				QuestionCount: 2,
				AnswerCount:   3,
				NSCount:       4,
				ARCount:       5,
			},
		},
		{
			expect: []byte{
				0b00000000, 0b00000001,
				0b10010000, 0b00000011,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			input: parser.Header{
				ID:            1,
				IsQuery:       false,
				ResponseCode:  parser.NAME_ERROR,
				OPCODE:        parser.STATUS,
				QuestionCount: 2,
				AnswerCount:   3,
				NSCount:       4,
				ARCount:       5,
			},
		},
		{
			expect: []byte{
				0b00000000, 0b00000001,
				0b10010000, 0b00000100,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			input: parser.Header{
				ID:            1,
				IsQuery:       false,
				OPCODE:        parser.STATUS,
				ResponseCode:  parser.NOT_IMPLEMENTED,
				QuestionCount: 2,
				AnswerCount:   3,
				NSCount:       4,
				ARCount:       5,
			},
		},
		{
			expect: []byte{
				0b00000000, 0b00000001,
				0b10010000, 0b00000101,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			input: parser.Header{
				ID:            1,
				IsQuery:       false,
				OPCODE:        parser.STATUS,
				ResponseCode:  parser.REFUSED,
				QuestionCount: 2,
				AnswerCount:   3,
				NSCount:       4,
				ARCount:       5,
			},
		},
	}
	for _, test := range tests {
		is, err := test.input.ToBinary()
		if err != nil {
			t.Fatalf("error should not be present")
		}
		CompareBytes(t, is, test.expect)
	}
}

func CompareBytes(t *testing.T, is []byte, expect []byte) {
	t.Helper()
	if string(is) != string(expect) {
		t.Logf("is:")
		for _, n := range is {
			t.Logf("%08b ", n) // prints 00000000 11111101
		}
		t.Logf("want:")
		for _, n := range expect {
			t.Logf("%08b ", n) // prints 00000000 11111101
		}
		t.Fatalf("strings dont match")
	}
}

func compareQuestion(t *testing.T, is parser.Question, expect parser.Question) {
	t.Helper()
	if is.Class != expect.Class {
		t.Fatalf("class dont match")
	}
	if is.Type != expect.Type {
		t.Fatalf("type dont match")
	}
	if len(is.Labels) != len(expect.Labels) {
		t.Fatalf("label length dont match")
	}
	for i := 0; i < len(is.Labels); i++ {
		if is.Labels[i] != expect.Labels[i] {
			t.Fatalf("label not match is: %s should: %s", is.Labels[i], expect.Labels[i])
		}
	}
}

func TestQuestionToBinary(t *testing.T) {
	tests := []struct {
		input  parser.Question
		expect []byte
	}{
		{
			input: parser.Question{
				Labels: []string{
					"blog",
					"example",
					"com",
				},
				Class: parser.IN,
			},
			expect: []byte{
				4, 'b', 'l', 'o', 'g',
				7, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
				3, 'c', 'o', 'm',
				0,

				0b00000000,
				0b00000000,

				0b00000000,
				0b00000001,
			},
		},
	}
	for _, test := range tests {
		is, err := test.input.ToBinary()
		if err != nil {
			t.Fatalf("should not error")
		}
		CompareBytes(t, is, test.expect)
	}
}

func compareAnswer(t *testing.T, is parser.Answer, expect parser.Answer) {
	t.Helper()
	if len(is.Labels) != len(expect.Labels) {
		t.Fatalf("labels length dont match up")
	}
	for i := 0; i < len(is.Labels); i++ {
		if is.Labels[i] != expect.Labels[i] {
			t.Fatalf("label dont match up is: %s want: %s", is.Labels[i], expect.Labels[i])
		}

	}
	if is.Type != expect.Type {
		t.Fatalf("type dont match up is: %d want: %d", is.Type, expect.Type)
	}
	if is.Class != expect.Class {
		t.Fatalf("class dont match up is: %d want: %d", is.Class, expect.Class)
	}
	if is.TTL != expect.TTL {
		t.Fatalf("ttl dont match up is: %d want: %d", is.TTL, expect.TTL)
	}

	if len(is.Data) != len(expect.Data) {
		t.Fatalf("data length dont match up is: %d want: %d", len(is.Data), len(expect.Data))
	}
	for i := 0; i < len(is.Data); i++ {
		if is.Data[i] != expect.Data[i] {
			t.Fatalf("data dont match is: %d want: %d i: %d", is.Data[i], expect.Data[i], i)

		}

	}
}

func TestAnswerToBinary(t *testing.T) {
	tests := []struct {
		input  parser.Answer
		expect []byte
	}{
		{
			input: parser.Answer{
				Labels: []string{
					"example",
					"com",
				},
				Type:  parser.A,
				TTL:   1024,
				Class: parser.IN,
				Data:  []byte{1, 1, 1, 1},
			},
			expect: []byte{
				7, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
				3, 'c', 'o', 'm',
				0,

				0b00000000, 0b00000001,

				0b00000000, 0b00000001,

				0b00000000, 0b00000000, 0b00000100, 0b00000000,

				0, 4, 1, 1, 1, 1,
			},
		},
	}
	for _, test := range tests {
		is, err := test.input.ToBinary()
		if err != nil {
			t.Fatalf("should not error")
		}
		CompareBytes(t, is, test.expect)
	}

}
