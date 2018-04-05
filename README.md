*documentation currently out of date after change of package structure and flags*


# Redirect
_redirect_ is a pure 'redirecting server', i.e. all requests will receive a http redirect status. _redirect_ uses hostname and request URL to determine the redirect target. It can be used as a backend for URL shorteners or domain forwarders.



## Use as standalone server
// TODO

### Setup of _redirect_


Testing of this setup is easiest with curl, e.g.
```
    # Start server on port 8080 (all interfaces) with admin interface on localhost
    ./redirect -admin localhost -ignoreError &

    # Add a redirect via REST API
    curl "http://localhost:8080/redirects/add?host=www.example.com&url=/&target=http://www.google.com"

    # Test redirect - should show redirect to www.google.com
    curl --header 'Host: www.example.org' http://127.0.0.1:8080
```

### Setup of DNS


## REST API
// TODO

## Config file

The config file is JSON for a map structure
```
{
 "Hosts": {
  "hostname": {
   "url": "target",
   "url": "target"
  },
  "hostname": {
   "url": "target"
  }
 }
}
```


An example for a config looks like follows
```
{
 "Hosts": {
  "host1.example.com": {
   "/": "http://google.com"
  }
 }
}
```

This is comparable to nginx configuration like follows (not showing listen address setup)
```
server {
	server_name host1.example.com;
	rewrite ^/(.*)$ http://google.com redirect;
}
```

A configuration file can also be created by starting the server with the `ignoreError` option, which will create an empty map at launch and save the active configuration when ending the server. Entries can e.g. be populated with the `adminclient` or any other use of the REST API.


## Use as library

The package `redirectserver` can be used as library. A minimal implementation for launching a redirect server looks like this:
```
package main

import (
	"log"

	"github.com/flo80/redirect/redirectserver"
)

func main() {
	server := redirectserver.NewServer(":8080") 
	log.Fatal(server.StartServer())
}
``` 
To start with the REST API the server should instead be defined as 
```	
server := redirectserver.NewServer(":8080", 
              redirectserver.WithAdmin("localhost"))
```

The redirector (i.e. the component saving the redirect data) could be replaced with any structure implementing the `Redirector` interface.
