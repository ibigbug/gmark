package gmark

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func removeS(i string) string {
	reg := regexp.MustCompile("\\s")
	return reg.ReplaceAllString(i, "")
}
func readFile(fp string) string {
	fd, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestConvert(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "headers",
			args: struct{ text string }{
				text: readFile("testdata/normal/headers.text"),
			},
			want: readFile("testdata/normal/headers.html"),
		},
		{
			name: "ordered_and_unordered_lists",
			args: struct{ text string }{
				text: readFile("testdata/normal/ordered_and_unordered_lists.text"),
			},
			want: readFile("testdata/normal/ordered_and_unordered_lists.html"),
		},
	}
	for i, tt := range tests {
		if i == 1 {
			continue
		}
		if got := Convert(tt.args.text); removeS(got) != removeS(tt.want) {
			t.Errorf("Convert() = %v, want %v", removeS(got), removeS(tt.want))
		}

	}
}
