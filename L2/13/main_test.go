package main

import (
	"reflect"
	"testing"
)

func TestParseFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		spec    string
		want    []int
		wantErr bool
	}{
		{
			name: "single and range",
			spec: "1,3-5",
			want: []int{1, 3, 4, 5},
		},
		{
			name: "duplicates are deduped",
			spec: "2,2,1-3",
			want: []int{1, 2, 3},
		},
		{
			name:    "invalid zero",
			spec:    "0,2",
			wantErr: true,
		},
		{
			name:    "invalid descending range",
			spec:    "4-2",
			wantErr: true,
		},
		{
			name:    "invalid token",
			spec:    "1,a",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseFields(tt.spec)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for spec %q", tt.spec)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for spec %q: %v", tt.spec, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseFields(%q) = %v, want %v", tt.spec, got, tt.want)
			}
		})
	}
}

func TestSelectFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		line      string
		delimiter string
		fields    []int
		separated bool
		want      string
		wantOK    bool
	}{
		{
			name:      "basic selection",
			line:      "a\tb\tc\td",
			delimiter: "\t",
			fields:    []int{1, 3},
			separated: false,
			want:      "a\tc",
			wantOK:    true,
		},
		{
			name:      "out of range ignored",
			line:      "a\tb",
			delimiter: "\t",
			fields:    []int{2, 4},
			separated: false,
			want:      "b",
			wantOK:    true,
		},
		{
			name:      "line without delimiter returned when no -s",
			line:      "plain",
			delimiter: "\t",
			fields:    []int{1, 2},
			separated: false,
			want:      "plain",
			wantOK:    true,
		},
		{
			name:      "line without delimiter skipped with -s",
			line:      "plain",
			delimiter: "\t",
			fields:    []int{1, 2},
			separated: true,
			want:      "",
			wantOK:    false,
		},
		{
			name:      "custom delimiter",
			line:      "x,y,z",
			delimiter: ",",
			fields:    []int{2, 3},
			separated: false,
			want:      "y,z",
			wantOK:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := selectFields(tt.line, tt.delimiter, tt.fields, tt.separated)
			if ok != tt.wantOK {
				t.Fatalf("selectFields() ok = %v, want %v", ok, tt.wantOK)
			}
			if got != tt.want {
				t.Fatalf("selectFields() got %q, want %q", got, tt.want)
			}
		})
	}
}
