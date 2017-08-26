package client

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/openconfig/gnmi/proto/gnmi"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		query         string
		parsedQueries []string
	}{
		{
			query:         "/a/b/c/d",
			parsedQueries: []string{"/a/b/c/d"},
		},
		{
			query:         "a/b/c/d,c/d/e/f",
			parsedQueries: []string{"a/b/c/d", "c/d/e/f"},
		},
		{
			query:         "/a/b/c/d[12/3=4]/e",
			parsedQueries: []string{"/a/b/c/d[12/3=4]/e"},
		},
	}
	for _, tt := range tests {
		got := ParseQuery(tt.query)
		if diff := pretty.Compare(tt.parsedQueries, got); diff != "" {
			t.Errorf("ParseQuery(%s) returned diff (-want +got):\n%s", tt.query, diff)
		}
	}
}

func TestToGetRequest(t *testing.T) {
	tests := []struct {
		queries    []string
		getRequest *gnmi.GetRequest
	}{
		{
			queries: []string{"/", "/a/b/c/d[a=123]/e", "c/d[\"a/b\"=\"12/3\"]"},
			getRequest: &gnmi.GetRequest{
				Path: []*gnmi.Path{
					{
						Elem: []*gnmi.PathElem{},
					},
					{
						Elem: []*gnmi.PathElem{
							{Name: "a"},
							{Name: "b"},
							{Name: "c"},
							{Name: "d", Key: map[string]string{"a": "123"}},
							{Name: "e"},
						},
					},
					{
						Elem: []*gnmi.PathElem{
							{Name: "c"},
							{Name: "d", Key: map[string]string{"a/b": "12/3"}},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got, err := ToGetRequest(tt.queries)
		if err != nil {
			t.Errorf("ToGetRequest(%s) returned error: %s", tt.queries, err)
		}
		if diff := pretty.Compare(tt.getRequest, got); diff != "" {
			t.Errorf("ToGetRequest(%s) returned diff (-want +got):\n%s", tt.queries, diff)
		}
	}
}
