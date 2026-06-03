package validate

import "testing"

type payload struct {
	Name string `json:"name" validate:"required"`
}

func TestStruct(t *testing.T) {
	if err := Struct(payload{Name: "handygo"}); err != nil {
		t.Fatalf("validate struct: %v", err)
	}
	if err := Struct(payload{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestFormat(t *testing.T) {
	err := Struct(payload{})
	fields := Format(err)
	if len(fields) != 1 {
		t.Fatalf("fields len = %d", len(fields))
	}
	if fields[0].Field != "name" {
		t.Fatalf("field = %q", fields[0].Field)
	}
	if fields[0].Tag != "required" {
		t.Fatalf("tag = %q", fields[0].Tag)
	}
}
