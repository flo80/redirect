package redirectserver

import (
	"encoding/json"
	"fmt"
	"log"
)

func _debug(format string, v ...interface{}) {
	log.Printf("+DEBUG+ "+format, v...)
}

// map[hostname][url]->redirect
type mapRedirect struct {
	Hosts map[string]map[string]string
}

func (red *mapRedirect) GetAllRedirects() map[string]map[string]string {
	return red.Hosts
}

func (red *mapRedirect) GetRedirectsForHost(hostname string) map[string]string {
	_debug("requested redirects for hostname %v", hostname)

	redirectHost, okHost := red.Hosts[hostname]
	if !okHost {
		return make(map[string]string)
	}

	return redirectHost
}

func (red *mapRedirect) GetRedirect(hostname, url string) string {
	_debug("requested redirects for hostname %v url%v", hostname, url)

	redirectHost, okHost := red.Hosts[hostname]
	if !okHost {
		return ""
	}

	redirect, okURL := redirectHost[url]
	if !okURL {
		return ""
	}

	return redirect
}

//GetTarget gets a redirect target for a host and URL
func (red *mapRedirect) GetTarget(hostname string, url string) (string, error) {
	_debug("GetTarget call for %v %v", hostname, url)

	target := red.GetRedirect(hostname, url)

	if target == "" {
		return target, fmt.Errorf("no redirect foud for %v%v", hostname, url)
	}

	_debug("redirect found for host %v and url %v, target %v", hostname, url, target)
	return target, nil
}

// AddRedirect adds or changes a new host and/or URL to the redirections.
func (red *mapRedirect) AddRedirect(hostname string, url string, target string) error {
	log.Printf("adding new entry %v%v -> %v", hostname, url, target)

	if red.Hosts == nil {
		_debug("creating new mapRedirect.Hosts map in AddRedirect")
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
func (red *mapRedirect) RemoveRedirectHost(hostname string) {
	if red.Hosts != nil {
		delete(red.Hosts, hostname)
	}
}

// RemoveURL deletes all existing redirections for a host
func (red *mapRedirect) RemoveRedirect(hostname, URL string) {
	if red.Hosts == nil {
		return
	}

	_, exists := red.Hosts[hostname]

	if !exists {
		return
	}

	delete(red.Hosts[hostname], URL)
}

//GetJSON of all redirects
func (red *mapRedirect) GetJSON() ([]byte, error) {
	return json.MarshalIndent(red, "", " ")
}

//SetJSON for all redirects
func (red *mapRedirect) SetJSON(b []byte) error {
	return json.Unmarshal(b, red)
}
