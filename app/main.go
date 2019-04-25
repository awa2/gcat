package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type JSON map[string]interface{}
type ARRAY []interface{}

func main() {
	var GcatData JSON = JSON{
		// "share":        GetWmicShare(),
		"Systeminfo":   GetUsers(),
		"MacAddresses": GetMacAddresses(),
	}

	bytes, err := json.MarshalIndent(GcatData, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
	fpw, err := os.Create("gcat.json")
	if err != nil {
		log.Fatal(err)
	}
	fpw.Write(bytes)
}

type User struct {
	Id   string `csv:"client_id"`
	Name string `csv:"client_name"`
	Age  string `csv:"client_age"`
}

//どうやらCSVの中身に空があるとエラーになるっぽいです。

// func GetWmicShare() JSON {
// 	// users := []*User{}
// 	// ここでCSV形式の文字列を受け取るコマンドを実行する。
// 	// arg := "list full /format:csv"
// 	cmd := exec.Command("wmic", "share", "list", "/format:csv")
// 	out, err := cmd.Output()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	data, err := conv(string(out))
// 	fmt.Println(string(data))
// 	// 	var in = `client_id,client_name,client_age
// 	// user01,"J, Smith", 21`

// 	r := csv.NewReader(strings.NewReader(data))
// 	records, err := r.ReadAll()
// 	fmt.Println(string("test"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	json := JSON{}

// 	var headers []string
// 	for i, record := range records {

// 		if i == 0 {
// 			headers = record
// 		} else {
// 			for index, header := range headers {
// 				json[header] = record[index]
// 			}
// 		}
// 	}
// 	return json
// }

func GetUsers() JSON {
	// users := []*User{}
	// ここでCSV形式の文字列を受け取るコマンドを実行する。
	cmd := exec.Command("systeminfo", "/fo", "csv")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	data, err := conv(string(out))
	// 	var in = `client_id,client_name,client_age
	// user01,"J, Smith", 21`

	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	json := JSON{}
	var headers []string
	for i, record := range records {
		if i == 0 {
			headers = record
		} else {
			for index, header := range headers {
				json[header] = record[index]
			}
		}
	}
	return json
}

func GetMacAddresses() []string {
	cmd := exec.Command("getmac")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	data, err := conv(string(out))
	r := regexp.MustCompile(`(..-){5}..`) // [0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}
	results := r.FindAllString(data, -1)
	return results
}

func conv(str string) (string, error) {
	strReader := strings.NewReader(str)
	decodedReader := transform.NewReader(strReader, japanese.ShiftJIS.NewDecoder())
	decoded, err := ioutil.ReadAll(decodedReader)
	if err != nil {
		return "", err
	}
	return string(decoded), err
}
