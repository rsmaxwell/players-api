package debug

import (
	"testing"
)

var (
	pkg              = NewPackage("debug")
	functionTestDump = NewFunction(pkg, "TestDump")
	functionFoo      = NewFunction(pkg, "foo")
	functionBar      = NewFunction(pkg, "bar")
)

func TestDump(t *testing.T) {
	f := functionTestDump

	f.DebugVerbose("name: %s", "hello")

	foo()
}

func foo() {
	f := functionFoo

	f.DebugVerbose("name: %s", "one")

	bar()
}

func bar() {
	f := functionBar

	f.DebugVerbose("name: %s", "two")
	f.Dump("dump at %s", "world")
}
