package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{
			name: "encode default",
			args: nil,
			want: []string{"encode", "6GCRPR78+CV"},
		},
		{
			name: "encode len 12",
			args: []string{"-lat", "-1.286386", "-lon", "36.817223", "-len", "12"},
			want: []string{"len=12", "6GCRPR78+CV"},
		},
		{
			name: "decode",
			args: []string{"-decode", "6GCRPR78+CV"},
			want: []string{"lo=(", "hi=(", "center=(", "len=10"},
		},
		{
			name:    "decode invalid",
			args:    []string{"-decode", "8FWC2_45+G6"},
			wantErr: true,
		},
		{
			name: "shorten",
			args: []string{"-shorten"},
			want: []string{"shorten -> +CV"},
		},
		{
			name: "recover",
			args: []string{"-recover", "+CV"},
			want: []string{"recover +CV -> 6GCRPR78+CV"},
		},
		{
			name:    "bad lat",
			args:    []string{"-lat", "abc"},
			wantErr: true,
		},
		{
			name:    "extra args",
			args:    []string{"extra"},
			wantErr: true,
		},
		{
			name:    "unknown flag",
			args:    []string{"-nope"},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := run(&buf, tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("run(%v) succeeded, want error; output %q", tc.args, buf.String())
				}
				return
			}
			if err != nil {
				t.Fatalf("run(%v): %v", tc.args, err)
			}
			out := buf.String()
			for _, want := range tc.want {
				if !strings.Contains(out, want) {
					t.Errorf("run(%v) output %q missing %q", tc.args, out, want)
				}
			}
		})
	}
}
