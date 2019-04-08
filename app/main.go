package main

import (
	"encoding/json"
	"flag"
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
		"MacAddresses": GetMacAddresses(),
		"info":         GetSysteminfo(),
		"Hostname":     GetHostname(),
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

func GetSysteminfo() []string {
	cmd := exec.Command("systeminfo", "/fo", "csv")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	data, err := conv(string(out))
	tmp := csvtojson(string(data))
	r := regexp.MustCompile(`.+:.+`)
	results := r.FindAllString(tmp, -1)
	return results
}

func GetHostname() []string {
	cmd := exec.Command("hostname")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	data, err := conv(string(out))
	r := regexp.MustCompile(`.+`) // [0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}
	results := r.FindAllString(data, -1)
	return results
}

var delimiter *string = flag.String("delimiter", ",", "specify separator (e.g. \"\\t\")")
var lazyQuote *bool = flag.Bool("lazyQuote", true, "allow lazyQuote")

func csvtojson(str string) string {
	tmpA := strings.TrimLeft(str, `"`)
	tmpB := strings.TrimRight(tmpA, "\r")
	tmpC := strings.Split(tmpB, "\n")

	var results string
	var arrKey []string
	var arrVal []string

	for i, v := range tmpC {
		if i == 0 {
			arrKey = strings.Split(v, `","`)
		}
		if v != "" && i > 0 {
			v = strings.TrimLeft(v, `"`)
			for j, w := range arrKey {
				arrVal = strings.Split(v, `","`)
				results = results + arrKey[j] + `":"` + arrVal[j] + "\n"
				w = w
			}
		}
	}
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

// func csvtojson(r io.Reader, columns []string) ([]byte, error) {
// 	rows := make([]map[string]string, 0)
// 	csvReader := csv.NewReader(r)
// 	csvReader.TrimLeadingSpace = true
// 	for {
// 		record, err := csvReader.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		row := make(map[string]string)
// 		for i, n := range columns {
// 			row[n] = record[i]
// 		}
// 		rows = append(rows, row)
// 	}
// 	data, err := json.MarshalIndent(&rows, "", "  ")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }
