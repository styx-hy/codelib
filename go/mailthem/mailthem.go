// Copyright 2013 Yang Hong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Program mailthem implements a simple SMTP mail deliver tool
package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/gopass"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"strings"
)

// SInfo student information
type SInfo struct {
	snum string
	name string
	addr string
}

// getAuthInfo get smtp authentication username and password
// from command line using gopass
func getAuthInfo() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')

	password, passerr := gopass.GetPass("Enter password: ")
	return strings.Trim(username, "\n"), password, passerr
}

// readStudentInfo reads student email addresses and name strings
// from the contact file
func readStudentInfo() []SInfo {
	// read email addresses from a file
	fmt.Print("Directory with address list: ")
	addrListPath, err := bufio.NewReader(os.Stdin).ReadString('\n')
	addrListFile, err := os.Open(strings.Trim(addrListPath, "\n"))
	if err != nil {
		panic(err)
	}

	addrListBuf := bufio.NewReader(addrListFile)

	// 50 students for initial set
	var sl []SInfo = make([]SInfo, 0, 50)
	for {
		line, err := addrListBuf.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		entry := strings.Split(strings.Trim(line, "\n"), " ")

		// make sure the line indeed has 3 parts
		if len(entry) == 3 {
			sl = append(sl, SInfo{entry[0], entry[1], entry[2]})
		}
		// the last line will return io.EOF
		if err == io.EOF {
			break
		}
	}
	return sl
}

// makeMessage generates message contents for each student
func makeMessage(each SInfo, header map[string]string) (msg string) {
	var buf bytes.Buffer // A Buffer needs no initialization.

	for k, v := range header {
		buf.WriteString(k)
		buf.WriteString(": ")
		buf.WriteString(v)
		buf.WriteString("\r\n")
	}
	buf.Write([]byte(each.name))
	buf.Write([]byte(" 你好：\n"))
	buf.Write([]byte("Lab 1 的答辩安排在 "))
	buf.Write([]byte(" 10 月 16 日晚上 6 点"))
	buf.Write([]byte("，地点在软件大楼 3402，请做好准备。如有疑问请回复邮件，谢谢。"))

	buf.WriteString("\n\n----------\n")
	buf.WriteString("Yang Hong\n")
	buf.WriteString("hy.styx@gmail.com\n")
	buf.WriteString("Institute of Parallel and Distributed Systems")
	// msg = each.name + " 你好：\n" +

	msg, _ = buf.ReadString(byte(0))
	return msg
}

func main() {
	username, password, err := getAuthInfo()
	if err != nil {
		fmt.Println(err)
	}

	sl := readStudentInfo()
	fmt.Println(sl)
	// to avoid server recognition, use gmail as default
	auth := smtp.PlainAuth("", username, password, "smtp.gmail.com")

	// fill message with student lab grades
	for _, each := range sl {
		title := "Lab 1 答辩安排"
		header := make(map[string]string)
		header["From"] = username
		header["To"] = each.addr
		header["Subject"] = title
		header["MIME-Version"] = "1.0"
		header["Content-Type"] = "text/plain; charset=\"utf-8\""
		// header["Content-Transfer-Encoding"] = "base64"

		msgstr := []byte(makeMessage(each, header))

		err = smtp.SendMail("smtp.gmail.com:587", auth, username, []string{each.addr}, msgstr)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("send ok %s\n", each.addr)
		}
	}
}
