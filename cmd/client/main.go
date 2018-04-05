package main

import "github.com/flo80/redirect/cmd/client/cmd"

//Build version (GIT SHA)
var Build = "development"

func main() {
	cmd.Execute(Build)
}
