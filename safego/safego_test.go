package safego

import "testing"

func TestRunRecoversPanic(t *testing.T) {
	called := false
	Run(func() {
		panic("boom")
	}, WithPanicHandler(func(recovered any, stack []byte) {
		called = true
		if recovered != "boom" {
			t.Fatalf("panic = %v", recovered)
		}
		if len(stack) == 0 {
			t.Fatal("stack is empty")
		}
	}))
	if !called {
		t.Fatal("panic handler was not called")
	}
}
