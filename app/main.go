package main

import (
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

func GetMacAddresses() []string {
	cmd := exec.Command("ifconfig")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	data, err := conv(string(out))
	r := regexp.MustCompile(`(..:){5}..`) // [0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}
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