package building

import "testing"

func TestMakeTarget(t *testing.T) {
	checkMakeTarget(t, "All", "all")
	checkMakeTarget(t, "AllStuff", "all-stuff")
	checkMakeTarget(t, "HTTP", "http")
	checkMakeTarget(t, "AllHTTP", "all-http")
	// $$$$ MAT: enhance camel to kebab conversion
	// checkMakeTarget(t, "HTTPThings", "http-things")
	// checkMakeTarget(t, "AllHTTPThings", "all-http-things")
}

func checkMakeTarget(t *testing.T, name, expected string) {
	actual, _ := makeTarget(name, "")
	if actual != expected {
		t.Errorf("makeTarget failed: expected %s, got %s", expected, actual)
	}
}
