package building

import "testing"

func TestMakeTarget(t *testing.T) {
	checkMakeTarget(t, "TargetAll", "all")
	checkMakeTarget(t, "TargetAllStuff", "all-stuff")
	checkMakeTarget(t, "TargetHTTP", "http")
	checkMakeTarget(t, "TargetAllHTTP", "all-http")
	// $$$$ MAT: enhance camel to kebab conversion
	// checkMakeTarget(t, "TargetHTTPThings", "http-things")
	// checkMakeTarget(t, "TargetAllHTTPThings", "all-http-things")
}

func checkMakeTarget(t *testing.T, name, expected string) {
	actual, _ := makeTarget(name, "")
	if actual != expected {
		t.Errorf("makeTarget failed: expected %s, got %s", expected, actual)
	}
}
