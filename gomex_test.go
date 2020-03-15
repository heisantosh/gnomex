package main

import (
	"testing"
)

func TestShellVersion(t *testing.T) {
	g := newGnomex()
	v := findGnomeShellVersion()

	if g.gnomeShellVersion != v {
		t.Errorf("GNOME Shell version expected %v but got %v", v, g.gnomeShellVersion)
	}
}

func TestHTTPClient(t *testing.T) {
	g := newGnomex()

	if g.client == nil {
		t.Errorf("expected HTTP client but got nil")
	}
}
