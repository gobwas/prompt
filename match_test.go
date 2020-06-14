package prompt

import "testing"

func TestMatch(t *testing.T) {
	for _, test := range []struct {
		name string
		str  string
		pat  string
		exp  bool
	}{
		{
			str: "fat cat",
			pat: "fct",
			exp: true,
		},
		{
			str: "fat cat",
			pat: "fcct",
			exp: false,
		},
		{
			str: "cartwheel",
			pat: "twl",
			exp: true,
		},
		{
			str: "cartwheel",
			pat: "cart",
			exp: true,
		},
		{
			str: "cartwheel",
			pat: "cw",
			exp: true,
		},
		{
			str: "cartwheel",
			pat: "ee",
			exp: true,
		},
		{
			str: "cartwheel",
			pat: "art",
			exp: true,
		},
		{
			str: "cartwheel",
			pat: "eeel",
			exp: false,
		},
		{
			str: "cartwheel",
			pat: "dog",
			exp: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if act, exp := match(test.str, test.pat), test.exp; act != exp {
				t.Errorf(
					"match(%#q, %#q) = %t; want %t",
					test.str, test.pat,
					act, exp,
				)
			}
		})
	}
}
