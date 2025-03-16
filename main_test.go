package main

import (
	"bytes"
	"testing"

	"github.com/pascal-sochacki/dns/internal/parser"
)

func TestParseHeader(t *testing.T) {
	tests := []struct {
		input  []byte
		expect parser.Header
	}{
		{
			input: []byte{
				0b00000000, 0b00000001,
				0b00000000, 0b00000000,
				0b00000000, 0b00000001,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
			},
			expect: parser.Header{
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
			input: []byte{
				0b00000000, 0b00000001,
				0b10001000, 0b00000001,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			expect: parser.Header{
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
			input: []byte{
				0b00000000, 0b00000001,
				0b10010000, 0b00000010,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			expect: parser.Header{
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
			input: []byte{
				0b00000000, 0b00000001,
				0b10010000, 0b00000011,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			expect: parser.Header{
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
			input: []byte{
				0b00000000, 0b00000001,
				0b10010000, 0b00000100,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			expect: parser.Header{
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
			input: []byte{
				0b00000000, 0b00000001,
				0b10010000, 0b00000101,
				0b00000000, 0b00000010,
				0b00000000, 0b00000011,
				0b00000000, 0b00000100,
				0b00000000, 0b00000101,
			},
			expect: parser.Header{
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
		is := parser.ParseHeader(bytes.NewBuffer(test.input))

		if is.ID != test.expect.ID {
			t.Fatalf("id dont match is %d wanted %d", is.ID, test.expect.ID)
		}
		if is.QuestionCount != test.expect.QuestionCount {
			t.Fatalf("question count dont match is %d wanted %d", is.QuestionCount, test.expect.QuestionCount)
		}
		if is.AnswerCount != test.expect.AnswerCount {
			t.Fatalf("answer count dont match is %d wanted %d", is.AnswerCount, test.expect.AnswerCount)
		}
		if is.NSCount != test.expect.NSCount {
			t.Fatalf("ns count dont match is %d wanted %d", is.NSCount, test.expect.NSCount)
		}
		if is.ARCount != test.expect.ARCount {
			t.Fatalf("ar count dont match is %d wanted %d", is.ARCount, test.expect.ARCount)
		}
		if is.IsQuery != test.expect.IsQuery {
			t.Fatalf("is query dont match is %t wanted %t", is.IsQuery, test.expect.IsQuery)
		}
		if is.OPCODE != test.expect.OPCODE {
			t.Fatalf("op code dont match is %d wanted %d", is.OPCODE, test.expect.OPCODE)
		}
		if is.ResponseCode != test.expect.ResponseCode {
			t.Fatalf("response code dont match is %d wanted %d", is.ResponseCode, test.expect.ResponseCode)
		}
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

		if string(is) != string(test.expect) {
			t.Logf("is:")
			for _, n := range is {
				t.Logf("%08b ", n) // prints 00000000 11111101
			}
			t.Logf("want:")
			for _, n := range test.expect {
				t.Logf("%08b ", n) // prints 00000000 11111101
			}
			t.Fatalf("strings dont match")
		}
	}
}

func TestParseQuestion(t *testing.T) {
	tests := []struct {
		input  []byte
		expect parser.Question
	}{
		{
			input: []byte{
				4, 'b', 'l', 'o', 'g',
				7, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
				3, 'c', 'o', 'm',
				0,

				0b00000000,
				0b00000001,

				0b00000000,
				0b00000001,
			},
			expect: parser.Question{
				Class: parser.IN,
				Type:  parser.A,
				Labels: []string{
					"blog",
					"example",
					"com",
				},
			},
		},
	}
	for _, test := range tests {
		is := parser.ParseQuestion(bytes.NewBuffer(test.input))
		if is.Class != test.expect.Class {
			t.Fatalf("class dont match")
		}
		if is.Type != test.expect.Type {
			t.Fatalf("type dont match")
		}
		if len(is.Labels) != len(test.expect.Labels) {
			t.Fatalf("label length dont match")
		}
		for i := 0; i < len(is.Labels); i++ {
			if is.Labels[i] != test.expect.Labels[i] {
				t.Fatalf("label not match is: %s should: %s", is.Labels[i], test.expect.Labels[i])
			}

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
		if string(is) != string(test.expect) {
			t.Logf("is:")
			for _, n := range is {
				t.Logf("%08b ", n) // prints 00000000 11111101
			}
			t.Fatalf("dont match")

		}
	}
}

func TestParseAnswer(t *testing.T) {
	tests := []struct {
		input  []byte
		expect parser.Answer
	}{
		{
			input: []byte{
				7, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
				3, 'c', 'o', 'm',
				0,

				0b00000000, 0b00000011,

				0b00000000, 0b00000001,

				0b00000000, 0b00000000, 0b00000100, 0b00000000,

				0, 4, 1, 2, 3, 4,
			},
			expect: parser.Answer{
				Labels: []string{
					"example",
					"com",
				},
				Type:  parser.MD,
				TTL:   1024,
				Class: parser.IN,
				Data:  []byte{1, 2, 3, 4},
			},
		},
	}
	for _, test := range tests {
		is := parser.ParseAnswer(bytes.NewBuffer(test.input))
		if len(is.Labels) != len(test.expect.Labels) {
			t.Fatalf("labels length dont match up")
		}
		for i := 0; i < len(is.Labels); i++ {
			if is.Labels[i] != test.expect.Labels[i] {
				t.Fatalf("label dont match up is: %s want: %s", is.Labels[i], test.expect.Labels[i])
			}

		}
		if is.Type != test.expect.Type {
			t.Fatalf("type dont match up is: %d want: %d", is.Type, test.expect.Type)
		}
		if is.Class != test.expect.Class {
			t.Fatalf("class dont match up is: %d want: %d", is.Class, test.expect.Class)
		}
		if is.TTL != test.expect.TTL {
			t.Fatalf("ttl dont match up is: %d want: %d", is.TTL, test.expect.TTL)
		}

		if len(is.Data) != len(test.expect.Data) {
			t.Fatalf("data length dont match up is: %d want: %d", len(is.Data), len(test.expect.Data))
		}
		for i := 0; i < len(is.Data); i++ {
			if is.Data[i] != test.expect.Data[i] {
				t.Fatalf("data dont match is: %d want: %d i: %d", is.Data[i], test.expect.Data[i], i)

			}

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
		if string(is) != string(test.expect) {
			t.Logf("is:")
			for _, n := range is {
				t.Logf("%08b ", n) // prints 00000000 11111101
			}
			t.Fatalf("dont match")
		}

	}

}
