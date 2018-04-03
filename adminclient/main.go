package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

type redirect struct {
	Hostname string //hostname of the redirector
	URL      string //URL on the hostname
	Target   string //forwarding address
}

type responseStatus struct {
	Status  bool
	Message string
	Content []redirect
}

func main() {
	serverAddress := flag.String("server", "localhost:8080", "Address of admin interface")
	flag.Parse()

	function := "ping"

	if flag.NArg() > 0 {
		function = flag.Arg(0)
	}

	paramNames := []string{"host", "url", "target"}
	var params []parameter

	l := flag.NArg() - 1

	if l > 1 {
		if len(paramNames) < l {
			l = len(paramNames)
		}
		params = make([]parameter, l)

		for i := 0; i < l; i++ {
			params[i] = parameter{paramNames[i], flag.Arg(i + 1)}
		}
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
