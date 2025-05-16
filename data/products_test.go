package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "tea",
		Price: 1.00,
		SKU:   "asdfghj-ghj-hjkl",
	}

	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
}
