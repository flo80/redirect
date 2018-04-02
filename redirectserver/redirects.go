package redirectserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func _debug(format string, v ...interface{}) {
	log.Printf("+DEBUG+ "+format, v...)
}

// map[hostname][url]->redirect
type hostRedirects struct {
	Hosts map[string]map[string]string
}

// Handler for http.HandleFunc
func (red *hostRedirects) Handler(w http.ResponseWriter, r *http.Request) {
	target, err := red.GetTarget(r.Host, r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		_debug("no redirect found: %v", err)
		return
	}
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	_debug("request received for host %v and url %v, redirected to %v", r.Host, r.URL, target)

}

// AddRedirect adds or changes a new host and/or URL to the redirections.
func (red *hostRedirects) AddRedirect(hostname, url, target string) error {
	log.Printf("adding new entry %v%v -> %v", hostname, url, target)

	if red.Hosts == nil {
		_debug("creating new hostRedirects.Hosts map in AddRedirect")
		red.Hosts = make(map[string]map[string]string)
	}

	_, exists := red.Hosts[hostname]
	if !exists {
		_debug("creating new url map for host %v in AddRedirect", hostname)
		red.Hosts[hostname] = make(map[string]string)
	}

	red.Hosts[hostname][url] = target

	return nil
}

// RemoveHost deletes all existing redirections for a host
func (red *hostRedirects) RemoveHost(hostname string) {
	if red.Hosts != nil {
		delete(red.Hosts, hostname)
	}
}

// RemoveURL deletes all existing redirections for a host
func (red *hostRedirects) RemoveURL(hostname, URL string) {
	if red.Hosts == nil {
		return
	}

	_, exists := red.Hosts[hostname]

	if !exists {
		return
	}

	delete(red.Hosts[hostname], URL)
}

//GetTarget gets a target for a host and URL
func (red *hostRedirects) GetTarget(hostname string, url string) (string, error) {
	_debug("getTarget call for %v %v", hostname, url)

	redirectHost, okHost := red.Hosts[hostname]
	if !okHost {
		return "", fmt.Errorf("cannot find host %v", hostname)
	}
	_debug("found hostmap %v", redirectHost)

	target, okURL := redirectHost[url]
	if !okURL {
		return "", fmt.Errorf("cannot find url %v on host %v", url, hostname)
	}

	_debug("redirect found for host %v and url %v, target %v", hostname, url, target)
	return target, nil
}

//GetJSON of all redirects
func (red *hostRedirects) GetJSON() ([]byte, error) {
	return json.MarshalIndent(red, "", " ")
}

//SetJSON for all redirects
func (red *hostRedirects) SetJSON(b []byte) error {
	return json.Unmarshal(b, red)
}
