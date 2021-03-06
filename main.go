package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
	"bytes"
	"os"
)

const baseURL = "https://servicearizona.com/webapp/vehicle/plates/personalizedChoiceSearch.do?plateChoice=001&choice="
const availableText = "Plate is available"
const validchars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func isAvailable(plate string) bool {
	resp, err := http.Get(baseURL + plate)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		panic(resp.Status)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return strings.Index(string(b), availableText) > -1
}

func main() {
	buff := bytes.NewBuffer(nil)

	for _,f := range validchars {
		for _,s := range validchars {
			plate := string(f) + string(s)
			if isAvailable(plate) {
				plate += " AVAILABLE"
			} else {
				plate += " TAKEN"
			}
			buff.Write([]byte(plate + "\r\n"))
			fmt.Println(plate)
			time.Sleep(time.Millisecond * 50)
		}
	}

	f, err := os.OpenFile("./plates.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := buff.Bytes()
	err = f.Truncate(int64(len(b)))
	if err != nil {
		panic(err)
	}
	_, err = f.WriteAt(b, 0)
	if err != nil {
		panic(err)
	}
}