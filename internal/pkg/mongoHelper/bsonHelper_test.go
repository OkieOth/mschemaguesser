package mongoHelper

import (
	"testing"
)

func TestFirstUpperCase(t *testing.T) {
	r1 := firstUpperCase("aaaa")
	if r1 != "Aaaa" {
		t.Errorf("aaaa != Aaaa: %s", r1)
		return
	}

	r2 := firstUpperCase("Aaaa")
	if r2 != "Aaaa" {
		t.Errorf("Aaaa != Aaaa: %s", r2)
		return
	}
}
