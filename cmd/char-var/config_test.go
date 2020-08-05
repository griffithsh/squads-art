package main

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	test := []byte(`{"Color":"FF00FF"}`)

	tc := struct {
		Color rgb
	}{}

	err := json.Unmarshal(test, &tc)

	if err != nil {
		t.Fatal(err)
	}

	if tc.Color.R != 255 {
		t.Errorf("want red 255, got %d", tc.Color.R)
	}
	if tc.Color.G != 0 {
		t.Errorf("want green 0, got %d", tc.Color.G)
	}
	if tc.Color.B != 255 {
		t.Errorf("want blue 255, got %d", tc.Color.B)
	}
}

func TestMarshalJSON(t *testing.T) {
	tc := struct{ Color rgb }{
		Color: rgb{R: 127, G: 42, B: 18},
	}

	b, err := json.Marshal(&tc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := string(b)
	want := `{"Color":"7f2a12"}`
	if got != want {
		t.Fatalf("want \n\t%s\ngot \n\t%s", want, got)
	}
}
