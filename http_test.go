package peccary

import "testing"

func TestGetControllerNameFromUrl(t *testing.T) {

	controllerName, err := getControllerNameFromUrl("/api/hello/world", "/api")
	if err != nil {
		t.Fatal(err)
	}
	if controllerName != "hello/world" {
		t.Fatalf("Pattern not correct")
	}

}
