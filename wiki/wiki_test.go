package main

import "testing"

func TestPages(t *testing.T) {
	pages, err := Pages()
	if err != nil {
		t.Fail()
	}
	for _, p := range pages {
		t.Log(p)
	}
}
