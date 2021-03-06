package storage

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// MapRedirect saves redirects in a map in memory
// Per default it uses a ruslog default logger, this can be overwritten with NewMapRedirector(logger)
type MapRedirect struct {
	Hosts  map[string]map[string]string // map[hostname][url]redirect
	logger *log.Logger                  // default logger
}

// NewMapRedirect allows to set the logger on the storage
func NewMapRedirect(logger *log.Logger) MapRedirect {
	r := MapRedirect{
		Hosts:  nil,
		logger: logger,
	}
	return r
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

func (red *MapRedirect) GetAllRedirects() []Redirect {
	return convertMapToSlice(red.Hosts)
}

func (red *MapRedirect) GetRedirectsForHost(hostname string) []Redirect {
	log.Debugf("requested redirects for hostname %v", hostname)

	redirectHost, okHost := red.Hosts[hostname]
	if !okHost {
		return nil
	}

	m := map[string]map[string]string{hostname: redirectHost}

	return convertMapToSlice(m)
}

func (red *MapRedirect) GetRedirect(hostname, url string) []Redirect {
	log.Debugf("requested redirects for hostname %v url%v", hostname, url)

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
func (red *MapRedirect) GetTarget(hostname string, url string) (string, error) {
	log.Debugf("GetTarget call for %v %v", hostname, url)

	target := red.GetRedirect(hostname, url)

	if target == nil || len(target) < 1 {
		return "", fmt.Errorf("no redirect foud for %v%v", hostname, url)
	}

	log.Debugf("redirect found for host %v and url %v, target %v", hostname, url, target)
	return target[0].Target, nil
}

// AddRedirect adds or changes a new host and/or URL to the redirections.
func (red *MapRedirect) AddRedirect(redirect Redirect) error {
	log.Printf("adding new entry %v%v -> %v", redirect.Hostname, redirect.URL, redirect.Target)

	if red.Hosts == nil {
		log.Debugf("creating new MapRedirect.Hosts map in AddRedirect")
		red.Hosts = make(map[string]map[string]string)
	}

	_, exists := red.Hosts[redirect.Hostname]
	if !exists {
		log.Debugf("creating new url map for host %v in AddRedirect", redirect.Hostname)
		red.Hosts[redirect.Hostname] = make(map[string]string)
	}

	red.Hosts[redirect.Hostname][redirect.URL] = redirect.Target

	return nil
}

// RemoveAllRedirectsForHost deletes all existing redirections for a host
func (red *MapRedirect) RemoveAllRedirectsForHost(redirect Redirect) {
	if red.Hosts != nil {
		delete(red.Hosts, redirect.Hostname)
	}
}

// RemoveRedirect deletes all existing redirections for a host
func (red *MapRedirect) RemoveRedirect(redirect Redirect) {
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
func (red *MapRedirect) GetJSON() ([]byte, error) {
	return json.MarshalIndent(red.Hosts, "", " ")
}

//SetJSON for all redirects
func (red *MapRedirect) SetJSON(b []byte) error {
	return json.Unmarshal(b, red.Hosts)
}
