package giotoast

import (
	"image/color"
	"testing"
)

func TestType_String(t *testing.T) {
	tests := []struct{
		name     string
		tt       Type
		expected string
	}{
		{"neutral", TypeNeutral, "Neutral"},
		{"success", TypeSuccess, "Success"},
		{"error", TypeError, "Error"},
		{"warning", TypeWarning, "Warning"},
		{"info", TypeInfo, "Info"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual string = test.tt.String()
			if test.expected != actual {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestType_Color(t *testing.T) {
	tests := []struct{
		name     string
		tt       Type
		expected color.NRGBA
	}{
		{"neutral", TypeNeutral, ColorBackground},
		{"success", TypeSuccess, ColorSuccess},
		{"error", TypeError, ColorError},
		{"warning", TypeWarning, ColorWarning},
		{"info", TypeInfo, ColorInfo},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual = test.tt.Color()
			if test.expected != actual {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestType_Icon(t *testing.T) {
	if nil != TypeNeutral.Icon() {
		t.Error("expected nil icon for TypeNeutral")
	}

	if nil == TypeSuccess.Icon() {
		t.Error("expected non-nil icon for TypeSuccess")
	}

	if nil == TypeError.Icon() {
		t.Error("expected non-nil icon for TypeError")
	}

	if nil == TypeWarning.Icon() {
		t.Error("expected non-nil icon for TypeWarning")
	}

	if nil == TypeInfo.Icon() {
		t.Error("expected non-nil icon for TypeInfo")
	}
}
