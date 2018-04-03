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

func convertMapToSlice(m map[string]map[string]string) []Redirect {
	r := make([]Redirect, 0)
	for hostname, urls := range m {
		for url, target := range urls {
			redirect := Redirect{hostname, url, target}
			r = append(r, redirect)
		}
	}
	return r
}

func (red *mapRedirect) GetAllRedirects() []Redirect {
	return convertMapToSlice(red.Hosts)
}

func (red *mapRedirect) GetRedirectsForHost(hostname string) []Redirect {
	_debug("requested redirects for hostname %v", hostname)

	redirectHost, okHost := red.Hosts[hostname]
	if !okHost {
		return nil
	}

	m := map[string]map[string]string{hostname: redirectHost}

	return convertMapToSlice(m)
}

func (red *mapRedirect) GetRedirect(hostname, url string) []Redirect {
	_debug("requested redirects for hostname %v url%v", hostname, url)

	redirectHost, okHost := red.Hosts[hostname]
	if !okHost {
		return nil
	}

	target, okURL := redirectHost[url]
	if !okURL {
		return nil
	}

	redirect := Redirect{hostname, url, target}

	return []Redirect{redirect}
}

//GetTarget gets a redirect target for a host and URL
func (red *mapRedirect) GetTarget(hostname string, url string) (string, error) {
	_debug("GetTarget call for %v %v", hostname, url)

	target := red.GetRedirect(hostname, url)

	if target == nil || len(target) < 1 {
		return "", fmt.Errorf("no redirect foud for %v%v", hostname, url)
	}

	_debug("redirect found for host %v and url %v, target %v", hostname, url, target)
	return target[0].Target, nil
}

// AddRedirect adds or changes a new host and/or URL to the redirections.
func (red *mapRedirect) AddRedirect(redirect Redirect) error {
	log.Printf("adding new entry %v%v -> %v", redirect.Hostname, redirect.URL, redirect.Target)

	if red.Hosts == nil {
		_debug("creating new mapRedirect.Hosts map in AddRedirect")
		red.Hosts = make(map[string]map[string]string)
	}

	_, exists := red.Hosts[redirect.Hostname]
	if !exists {
		_debug("creating new url map for host %v in AddRedirect", redirect.Hostname)
		red.Hosts[redirect.Hostname] = make(map[string]string)
	}

	red.Hosts[redirect.Hostname][redirect.URL] = redirect.Target

	return nil
}

// RemoveHost deletes all existing redirections for a host
func (red *mapRedirect) RemoveAllRedirectsForHost(redirect Redirect) {
	if red.Hosts != nil {
		delete(red.Hosts, redirect.Hostname)
	}
}

// RemoveURL deletes all existing redirections for a host
func (red *mapRedirect) RemoveRedirect(redirect Redirect) {
	if red.Hosts == nil {
		return
	}

	_, exists := red.Hosts[redirect.Hostname]

	if !exists {
		return
	}

	delete(red.Hosts[redirect.Hostname], redirect.URL)
}

//GetJSON of all redirects
func (red *mapRedirect) GetJSON() ([]byte, error) {
	return json.MarshalIndent(red, "", " ")
}

//SetJSON for all redirects
func (red *mapRedirect) SetJSON(b []byte) error {
	return json.Unmarshal(b, red)
}
