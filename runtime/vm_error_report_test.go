package runtime_test

import "testing"

func TestVMErrorInfo(t *testing.T) {
	expectError(t, `a := 5
a + "boo"`,
		"Runtime Error: invalid operation: int + string\n\tat test:2:1")

	expectError(t, `a := 5
b := a(5)`,
		"Runtime Error: not callable: int\n\tat test:2:6")

	expectError(t, `a := 5
b := {}
b.x.y = 10`,
		"Runtime Error: not index-assignable: undefined\n\tat test:3:1")

	expectError(t, `
a := func() {
	b := 5
	b += "foo"
}
a()`,
		"Runtime Error: invalid operation: int + string\n\tat test:4:2")

	expectErrorWithUserModules(t, `a := 5
	a + import("mod1")`, map[string]string{
		"mod1": `export "foo"`,
	}, ": invalid operation: int + string\n\tat test:2:2")

	expectErrorWithUserModules(t, `a := import("mod1")()`, map[string]string{
		"mod1": `
export func() {
	b := 5
	return b + "foo"
}`,
	}, "Runtime Error: invalid operation: int + string\n\tat mod1:4:9")

	expectErrorWithUserModules(t, `a := import("mod1")()`, map[string]string{
		"mod1": `export import("mod2")()`,
		"mod2": `
export func() {
	b := 5
	return b + "foo"
}`,
	}, "Runtime Error: invalid operation: int + string\n\tat mod2:4:9")
}
