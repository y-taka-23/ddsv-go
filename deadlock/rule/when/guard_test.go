package when_test

import (
	"testing"

	"github.com/y-taka-23/ddsv-go/deadlock/rule/vars"
	"github.com/y-taka-23/ddsv-go/deadlock/rule/when"
)

func TestIs(t *testing.T) {

	tests := []struct {
		name      string
		var_      vars.Name
		val       int
		in        vars.Shared
		want      bool
		wantError bool
	}{
		{
			name: "declared as equal", var_: "x", val: 42,
			in:   vars.Shared{"x": 42},
			want: true, wantError: false,
		},
		{
			name: "declared as not equal", var_: "x", val: 42,
			in:   vars.Shared{"x": 1},
			want: false, wantError: false,
		},
		{
			name: "undeclared", var_: "x", val: 42,
			in:   vars.Shared{},
			want: false, wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := when.Var(tt.var_).Is(tt.val)(tt.in)
			if tt.wantError && err == nil {
				t.Fatalf("want error, but has no error")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("want no error, but has error %v", err)
			}
			if !tt.wantError && got != tt.want {
				t.Fatalf("want %+v, but %+v", tt.want, got)
			}
		})
	}

}

func TestIsNot(t *testing.T) {

	tests := []struct {
		name      string
		var_      vars.Name
		val       int
		in        vars.Shared
		want      bool
		wantError bool
	}{
		{
			name: "declared as equal", var_: "x", val: 42,
			in:   vars.Shared{"x": 42},
			want: false, wantError: false,
		},
		{
			name: "declared as not equal", var_: "x", val: 42,
			in:   vars.Shared{"x": 1},
			want: true, wantError: false,
		},
		{
			name: "undeclared", var_: "x", val: 42,
			in:   vars.Shared{},
			want: false, wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := when.Var(tt.var_).IsNot(tt.val)(tt.in)
			if tt.wantError && err == nil {
				t.Fatalf("want error, but has no error")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("want no error, but has error %v", err)
			}
			if !tt.wantError && got != tt.want {
				t.Fatalf("want %+v, but %+v", tt.want, got)
			}
		})
	}

}

func TestIsLessThan(t *testing.T) {

	tests := []struct {
		name      string
		var_      vars.Name
		val       int
		in        vars.Shared
		want      bool
		wantError bool
	}{
		{
			name: "declared as less", var_: "x", val: 42,
			in:   vars.Shared{"x": 41},
			want: true, wantError: false,
		},
		{
			name: "declared as equal", var_: "x", val: 42,
			in:   vars.Shared{"x": 42},
			want: false, wantError: false,
		},
		{
			name: "declared as greater", var_: "x", val: 42,
			in:   vars.Shared{"x": 43},
			want: false, wantError: false,
		},
		{
			name: "undeclared", var_: "x", val: 42,
			in:   vars.Shared{},
			want: false, wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := when.Var(tt.var_).IsLessThan(tt.val)(tt.in)
			if tt.wantError && err == nil {
				t.Fatalf("want error, but has no error")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("want no error, but has error %v", err)
			}
			if !tt.wantError && got != tt.want {
				t.Fatalf("want %+v, but %+v", tt.want, got)
			}
		})
	}

}

func TestIsGreaterThan(t *testing.T) {

	tests := []struct {
		name      string
		var_      vars.Name
		val       int
		in        vars.Shared
		want      bool
		wantError bool
	}{
		{
			name: "declared as less", var_: "x", val: 42,
			in:   vars.Shared{"x": 41},
			want: false, wantError: false,
		},
		{
			name: "declared as equal", var_: "x", val: 42,
			in:   vars.Shared{"x": 42},
			want: false, wantError: false,
		},
		{
			name: "declared as greater", var_: "x", val: 42,
			in:   vars.Shared{"x": 43},
			want: true, wantError: false,
		},
		{
			name: "undeclared", var_: "x", val: 42,
			in:   vars.Shared{},
			want: false, wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := when.Var(tt.var_).IsGreaterThan(tt.val)(tt.in)
			if tt.wantError && err == nil {
				t.Fatalf("want error, but has no error")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("want no error, but has error %v", err)
			}
			if !tt.wantError && got != tt.want {
				t.Fatalf("want %+v, but %+v", tt.want, got)
			}
		})
	}

}
