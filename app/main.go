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

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type JSON map[string]interface{}
type ARRAY []interface{}

func main() {
	var GcatData JSON = JSON{
		"Systeminfo":   GetUsers(),
		"MacAddresses": GetMacAddresses(),
		"Share":        GetWmicShare(),
		"Useraccount":  GetWmicUseraccount(),
		"QFfe":         GetWmicQfe(),
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
	// Aho  string `csv:"client_aho"`
}

type Share struct {
	Node   string `csv:"Node"`
	Name   string `csv:"Name"`
	Status string `csv:"Status"`
	Type   int    `csv:"Type"`
}

func GetWmicShare() []*JSON {
	Shares := []*Share{}
	// ここでCSV形式の文字列を受け取るコマンドを実行する。
	cmd := exec.Command("wmic", "share", "list", "/format:csv") // cmd := exec.Command("wmic", "timezone", "list", "brief", "/format:csv")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	data, err := conv(string(out))
	rep := regexp.MustCompile(`\nN`) //先頭行（空っぽ）を削除（置換）する
	data = rep.ReplaceAllString(data, "N")
	rep = regexp.MustCompile(`\r`) //\rを削除する
	data = rep.ReplaceAllString(data, "")
	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	jsons := []*JSON{}
	var headers []string
	for i, record := range records {
		if i == 0 {
			headers = record
		} else {
			json := JSON{}
			for index, header := range headers {
				json[header] = record[index]
			}
			jsons = append(jsons, &json)
		}
	}
	// fmt.Println(jsons)
	err = gocsv.UnmarshalString(data, &Shares)
	if err != nil {
		log.Fatal(err)
	}
	return jsons
}

func GetWmicUseraccount() []*JSON {
	Shares := []*Share{}
	// ここでCSV形式の文字列を受け取るコマンドを実行する。
	cmd := exec.Command("wmic", "useraccount", "list", "/format:csv") // cmd := exec.Command("wmic", "timezone", "list", "brief", "/format:csv")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	data, err := conv(string(out))
	rep := regexp.MustCompile(`\nN`) //先頭行（空っぽ）を削除（置換）する
	data = rep.ReplaceAllString(data, "N")
	rep = regexp.MustCompile(`\r`) //\rを削除する
	data = rep.ReplaceAllString(data, "")
	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	jsons := []*JSON{}
	var headers []string
	for i, record := range records {
		if i == 0 {
			headers = record
		} else {
			json := JSON{}
			for index, header := range headers {
				json[header] = record[index]
			}
			jsons = append(jsons, &json)
		}
	}
	// fmt.Println(jsons)
	err = gocsv.UnmarshalString(data, &Shares)
	if err != nil {
		log.Fatal(err)
	}
	return jsons
}

func GetWmicQfe() []*JSON {
	Shares := []*Share{}
	// ここでCSV形式の文字列を受け取るコマンドを実行する。
	cmd := exec.Command("wmic", "qfe", "list", "/format:csv") // cmd := exec.Command("wmic", "timezone", "list", "brief", "/format:csv")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	data, err := conv(string(out))

	data = Header_Del(data)
	// rep := regexp.MustCompile(`\nN`) //先頭行（空っぽ）を削除（置換）する
	// data = rep.ReplaceAllString(data, "N")
	// rep = regexp.MustCompile(`\r`) //\rを削除する
	// data = rep.ReplaceAllString(data, "")

	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	jsons := []*JSON{}
	var headers []string
	for i, record := range records {
		if i == 0 {
			headers = record
		} else {
			json := JSON{}
			for index, header := range headers {
				json[header] = record[index]
			}
			jsons = append(jsons, &json)
		}
	}
	// fmt.Println(jsons)
	err = gocsv.UnmarshalString(data, &Shares)
	if err != nil {
		log.Fatal(err)
	}
	return jsons
}

//↓お手本
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

func Header_Del(str string) {
	rep := regexp.MustCompile(`\nN`) //先頭行（空っぽ）を削除（置換）する
	data := rep.ReplaceAllString(str, "N")
	rep = regexp.MustCompile(`\r`) //\rを削除する
	data = rep.ReplaceAllString(data, "")

	return data, ""
}
