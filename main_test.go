package main

import "testing"

func TestParseHeader(t *testing.T) {
	tests := []struct {
		input  []byte
		expect Header
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
			expect: Header{
				ID:            1,
				IsQuery:       true,
				ResponseCode:  NO_ERROR,
				OPCODE:        QUERY,
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
			expect: Header{
				ID:            1,
				IsQuery:       false,
				ResponseCode:  FORMAT_ERROR,
				OPCODE:        IQUERY,
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
			expect: Header{
				ID:            1,
				IsQuery:       false,
				ResponseCode:  SERVER_FAILURE,
				OPCODE:        STATUS,
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
			expect: Header{
				ID:            1,
				IsQuery:       false,
				ResponseCode:  NAME_ERROR,
				OPCODE:        STATUS,
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
			expect: Header{
				ID:            1,
				IsQuery:       false,
				OPCODE:        STATUS,
				ResponseCode:  NOT_IMPLEMENTED,
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
			expect: Header{
				ID:            1,
				IsQuery:       false,
				OPCODE:        STATUS,
				ResponseCode:  REFUSED,
				QuestionCount: 2,
				AnswerCount:   3,
				NSCount:       4,
				ARCount:       5,
			},
		},
	}
	for _, test := range tests {
		is := ParseHeader(test.input)

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
		input  Header
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
			input: Header{
				ID:            1,
				IsQuery:       true,
				ResponseCode:  NO_ERROR,
				OPCODE:        QUERY,
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
			input: Header{
				ID:            1,
				IsQuery:       false,
				ResponseCode:  FORMAT_ERROR,
				OPCODE:        IQUERY,
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
			input: Header{
				ID:            1,
				IsQuery:       false,
				ResponseCode:  SERVER_FAILURE,
				OPCODE:        STATUS,
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
			input: Header{
				ID:            1,
				IsQuery:       false,
				ResponseCode:  NAME_ERROR,
				OPCODE:        STATUS,
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
			input: Header{
				ID:            1,
				IsQuery:       false,
				OPCODE:        STATUS,
				ResponseCode:  NOT_IMPLEMENTED,
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
			input: Header{
				ID:            1,
				IsQuery:       false,
				OPCODE:        STATUS,
				ResponseCode:  REFUSED,
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
		expect Question
	}{
		{
			input: []byte{
				0b00000100, 0b01100010, 0b01101100, 0b01101111, 0b01100111,
				0b00000111, 0b01100101, 0b01111000, 0b01100001, 0b01101101, 0b01110000, 0b01101100, 0b01100101,
				0b00000011, 0b01100011, 0b01101111, 0b01101101,

				0b00000000,

				0b00000000,
				0b00000001,

				0b00000000,
				0b00000001,
			},
			expect: Question{
				Class: IN,
				Type:  A,
				Labels: []string{
					"blog",
					"example",
					"com",
				},
			},
		},
	}
	for _, test := range tests {
		is := ParseQuestion(test.input)
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

func TestToBinary(t *testing.T) {
	tests := []struct {
		input  Question
		expect []byte
	}{
		{
			input: Question{
				Labels: []string{
					"blog",
					"example",
					"com",
				},
				Class: IN,
			},
			expect: []byte{
				0b00000100, 0b01100010, 0b01101100, 0b01101111, 0b01100111,
				0b00000111, 0b01100101, 0b01111000, 0b01100001, 0b01101101, 0b01110000, 0b01101100, 0b01100101,
				0b00000011, 0b01100011, 0b01101111, 0b01101101,

				0b00000000,

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
