package main

import (
	"bytes"
	"log/slog"
	"net"
	"os"

	"github.com/pascal-sochacki/dns/internal/parser"
)

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
		responseHeader := parser.Header{
			ID:            requestheader.ID,
			QuestionCount: 1,
			AnswerCount:   1,
		}
		request, _ = responseHeader.ToBinary()
		response.Write(request)

		question, _ := resquestquestion.ToBinary()
		response.Write(question)

		answer := parser.Answer{
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

func ParseMessage(buf []byte) (parser.Header, parser.Question) {
	header := parser.ParseHeader(bytes.NewBuffer(buf[:4*3]))
	question := parser.ParseQuestion(bytes.NewBuffer(buf[4*3:]))
	slog.Info("question", "type", question.Type)
	return header, question
}
