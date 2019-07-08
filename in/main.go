package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	indata, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	output, err := Execute(indata)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(output)
}

//Execute - provides in capability
func Execute(input []byte) (string, error) {
	var outdata struct {
		Version interface{} `json:"version"`
	}

	err := json.Unmarshal(input, &outdata)
	if err != nil {
		return "", err
	}
	if outdata.Version == nil {
		return "", errors.New("missing version")
	}
	outbytes, err := json.Marshal(outdata)
	return string(outbytes), err
}
