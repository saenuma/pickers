package main

import (
	_ "embed"
)

//go:embed Roboto-Light.ttf
var DefaultFont []byte

//go:embed colors.txt
var Colors2 []byte
