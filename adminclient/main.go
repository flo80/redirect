package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

type responseStatus struct {
	Status  bool
	Message string
	Content interface{}
}

func main() {
	serverAddress := flag.String("server", "localhost:8080", "Address of admin interface")
	flag.Parse()

	function := "ping"

	if flag.NArg() > 0 {
		function = flag.Arg(0)
	}

	paramNames := []string{"host", "url", "target"}
	params := make([]parameter, flag.NArg()-1)

	l := flag.NArg() - 1
	if len(paramNames) < l {
		l = len(paramNames)
	}

	for i := 0; i < l; i++ {
		params[i] = parameter{paramNames[i], flag.Arg(i + 1)}
	}

	response, err := requestFromServer(*serverAddress, function, params)

	if err != nil {
		fmt.Printf("error in request: %v\n", err)
	}

	fmt.Printf("response received: %v\n", response)
}

type parameter struct {
	key   string
	value string
}

func requestFromServer(serverAddress string, function string, params []parameter) (*responseStatus, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/redirects/%v", serverAddress, function), nil)
	if err != nil {
		return nil, fmt.Errorf("could not build request: %v", err)
	}

	if params != nil {
		q := req.URL.Query()
		for _, param := range params {
			q.Add(param.key, param.value)
		}
		req.URL.RawQuery = q.Encode()

		fmt.Printf("DEBUG query %v\n", req.URL.RawQuery)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found: %v", req.URL)

	}

	var response responseStatus

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}

	return &response, nil

}
