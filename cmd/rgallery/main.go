package main

import "github.com/robbymilo/rgallery/pkg/rgallery"

var Commit string
var Tag string

func main() {
	rgallery.SetupApp(Commit, Tag)
}
