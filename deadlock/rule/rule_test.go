package rule_test

import (
	"errors"
	"testing"

	"github.com/y-taka-23/ddsv-go/deadlock/rule"
)

func TestCopy(t *testing.T) {

	tests := []struct {
		name      string
		from      rule.VarName
		to        rule.VarName
		in        rule.SharedVars
		want      rule.SharedVars
		wantError bool
	}{
		{
			name: "defined to defined", from: "x", to: "y",
			in:        rule.SharedVars{"x": 42, "y": 1},
			want:      rule.SharedVars{"x": 42, "y": 42},
			wantError: false,
		},
		{
			name: "defined to undefined", from: "x", to: "y",
			in:        rule.SharedVars{"x": 42},
			want:      rule.SharedVars{},
			wantError: true,
		},
		{
			name: "undefined to defined", from: "x", to: "y",
			in:        rule.SharedVars{"y": 42},
			want:      rule.SharedVars{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.Copy(tt.from, tt.to)(tt.in)
			if tt.wantError && errors.Is(err, nil) {
				t.Fatalf("want error, but has no error")
			}
			if !tt.wantError && !errors.Is(err, nil) {
				t.Fatalf("want no error, but has error %v", err)
			}
			if !tt.wantError && !eqVars(got, tt.want) {
				t.Fatalf("want %+v, but %+v", tt.want, got)
			}
		})
	}

}

func TestSet(t *testing.T) {

	tests := []struct {
		name      string
		val       int
		to        rule.VarName
		in        rule.SharedVars
		want      rule.SharedVars
		wantError bool
	}{
		{
			name: "to defined", val: 42, to: "x",
			in:        rule.SharedVars{"x": 1},
			want:      rule.SharedVars{"x": 42},
			wantError: false,
		},
		{
			name: "to undefined", val: 42, to: "x",
			in:        rule.SharedVars{},
			want:      rule.SharedVars{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.Set(tt.val, tt.to)(tt.in)
			if tt.wantError && errors.Is(err, nil) {
				t.Fatalf("want error, but has no error")
			}
			if !tt.wantError && !errors.Is(err, nil) {
				t.Fatalf("want no error, but has error %v", err)
			}
			if !tt.wantError && !eqVars(got, tt.want) {
				t.Fatalf("want %+v, but %+v", tt.want, got)
			}
		})
	}

}

func TestAdd(t *testing.T) {

	tests := []struct {
		name      string
		val       int
		to        rule.VarName
		in        rule.SharedVars
		want      rule.SharedVars
		wantError bool
	}{
		{
			name: "to defined", val: 42, to: "x",
			in:        rule.SharedVars{"x": 1},
			want:      rule.SharedVars{"x": 43},
			wantError: false,
		},
		{
			name: "to undefined", val: 42, to: "x",
			in:        rule.SharedVars{},
			want:      rule.SharedVars{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rule.Add(tt.val, tt.to)(tt.in)
			if tt.wantError && errors.Is(err, nil) {
				t.Fatalf("want error, but has no error")
			}
			if !tt.wantError && !errors.Is(err, nil) {
				t.Fatalf("want no error, but has error %v", err)
			}
			if !tt.wantError && !eqVars(got, tt.want) {
				t.Fatalf("want %+v, but %+v", tt.want, got)
			}
		})
	}

}

func eqVars(got, want rule.SharedVars) bool {
	if len(got) != len(want) {
		return false
	}
	for x, _ := range got {
		if got[x] != want[x] {
			return false
		}
	}
	return true
}
