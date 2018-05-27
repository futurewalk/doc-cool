package main

import (
	"testing"
	"doc-cool/cool"
)

func TestStr(t *testing.T)  {
	cool.GetIgnoreFile("push.proto,login.proto")
}

func TestExport(t *testing.T)  {
	cool.Export()
}