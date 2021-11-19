package utils

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true item contained in list",
			args: args{
				s:   []string{"foo", "bar", "baz"},
				str: "foo",
			},
			want: true,
		},
		{
			name: "should return false if not contained in list",
			args: args{
				s:   []string{"foo", "bar", "baz"},
				str: "wakka",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.str); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type args struct {
		source []string
		match  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "should find match",
			args: args{
				source: []string{"foo", "bar", "baz"},
				match:  "foo",
			},
			want: []string{"foo"},
		},
		{
			name: "should return nothing if no match",
			args: args{
				source: []string{"foo", "bar", "baz"},
				match:  "wakka",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Find(tt.args.source, tt.args.match); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
