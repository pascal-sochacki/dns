package main

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"github.com/pascal-sochacki/dns/internal/parser"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", "192.203.230.10:53")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	header := parser.Header{
		ID:            12,
		IsQuery:       true,
		QuestionCount: 1,
	}
	question := parser.Question{
		Labels: []string{"eu"},
		Type:   parser.A,
		Class:  parser.IN,
	}
	request := new(bytes.Buffer)

	buf, _ := header.ToBinary()
	request.Write(buf)

	buf, _ = question.ToBinary()
	request.Write(buf)

	_, err = conn.Write(request.Bytes())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reponse := make([]byte, 1024)
	_, _, _ = conn.ReadFrom(reponse)
	buffer := bytes.NewBuffer(reponse)
	responseHeader := parser.ParseHeader(buffer)
	println(responseHeader.String())
	for i := 0; i < int(responseHeader.QuestionCount); i++ {
		question := parser.ParseQuestion(buffer)
		println(question.String())
	}

	println("AnswerCount")
	for i := 0; i < int(responseHeader.AnswerCount); i++ {
		question := parser.ParseAnswer(buffer)
		println(question.String())
	}
	println("NSCount")
	for i := 0; i < int(responseHeader.NSCount); i++ {
		question := parser.ParseAnswer(buffer)
		println(question.String())
	}

}
