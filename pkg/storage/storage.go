package storage

//Redirect entry declaration
type Redirect struct {
	Hostname string //hostname of the redirector
	URL      string //URL on the hostname
	Target   string //forwarding address
}

// Redirector interface
type Redirector interface {
	GetAllRedirects() []Redirect                                      // Get all redirects known to redirects
	GetRedirectsForHost(hostname string) []Redirect                   // Get all redirects for a specific hostname
	GetRedirect(hostname string, url string) []Redirect               // Get redirect for a specific hostname & url (should be only one)
	AddRedirect(redirect Redirect) error                              // Add a new redirect for a hostname & url
	RemoveRedirect(redirect Redirect)                                 // Remove a redirect specific to hostname & url
	RemoveAllRedirectsForHost(redirect Redirect)                      // Remove all redirects for a hostname
	GetTarget(hostname string, url string) (target string, err error) // Return the redirect target for the hostname & url
}
