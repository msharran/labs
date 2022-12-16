package search

import (
	"bytes"
	"strings"
	"testing"
)

const CONTENT = `Lorem ipsum dolor sit amet,
consectetur adipiscing elit. Curabitur et magna
vel dolor porttitor ullamcorper vitae non lacus. 
Etiam consectetur, nibh quis placerat posuere,
leo eros lacinia massa, eu viverra tellus augue 
`

func TestSearch(t *testing.T) {
	testCases := []struct {
		name string
		args ExecArgs
		want string
	}{
		{
			"simple search",
			ExecArgs{
				Pattern: "quis",
			},
			"Etiam consectetur, nibh quis placerat posuere,\n",
		},
		{
			"search with line numbers",
			ExecArgs{
				Pattern:         "Lorem",
				ShowLineNumbers: true,
			},
			"1 Lorem ipsum dolor sit amet,\n",
		},
		{
			"search case insensitive",
			ExecArgs{
				Pattern:         "lorem",
				CaseInsensitive: true,
			},
			"Lorem ipsum dolor sit amet,\n",
		},
		{
			"search line count",
			ExecArgs{
				Pattern:       "et",
				ShowLineCount: true,
			},
			"3\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			w := &bytes.Buffer{}

			err := Exec(w, strings.NewReader(CONTENT), tc.args)

			if err != nil {
				t.Error(err)
			}

			if w.String() != tc.want {
				t.Errorf("got: %q want: %q", w.String(), tc.want)
			}
		})
	}
}
