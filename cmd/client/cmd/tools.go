package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type redirect struct {
	Hostname string //hostname of the redirector
	URL      string //URL on the hostname
	Target   string //target address
}

type response struct {
	Status  bool
	Message string
	Content []redirect
}

type parameter struct {
	key   string
	value string
}

func createParamsFromArgs(args []string) []parameter {
	paramNames := []string{"host", "url", "target"}
	var params []parameter

	l := len(args)

	if l > 0 {
		if len(paramNames) < l {
			l = len(paramNames)
		}
		params = make([]parameter, l)

		for i := 0; i < l; i++ {
			params[i] = parameter{paramNames[i], args[i]}
		}
	}

	return params
}

func requestFromServer(function string, args []string) error {
	server, err := rootCmd.PersistentFlags().GetString("server")
	if err != nil {
		return fmt.Errorf("Server address not available: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/redirects/%v", server, function), nil)
	if err != nil {
		return fmt.Errorf("could not build request: %v", err)
	}

	params := createParamsFromArgs(args)
	if params != nil {
		q := req.URL.Query()
		for _, param := range params {
			q.Add(param.key, param.value)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Server / API not found at %v \n\n", server)
		os.Exit(1)
	}

	var response response

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("could not decode response: %v", err)
	}

	return processResponse(&response)

}

func processResponse(response *response) error {

	if !response.Status {
		return fmt.Errorf("Operation was not sucessful on server, error %v", response.Message)
	}

	fmt.Printf("Operation successful (%v) \n\n", response.Message)

	if len(response.Content) > 0 {
		fmt.Printf("%-30s %-10s %-50s \n", "Hostname", "URL", "Target")
		fmt.Printf("%-30s %-10s %-50s \n", "--------", "---", "------")
		for _, r := range response.Content {
			fmt.Printf("%-30s %-10s %-50s \n", r.Hostname, r.URL, r.Target)
		}
		fmt.Println()
	}
	return nil
}
