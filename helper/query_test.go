package helper

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
			query:         "/a/b/c/d/",
			parsedQueries: []string{"a/b/c/d"},
		},
		{
			query:         "/a/b/c/d/,c/d/e/f",
			parsedQueries: []string{"a/b/c/d", "c/d/e/f"},
		},
		{
			query:         "/a/b/c/d[123]/e",
			parsedQueries: []string{"a/b/c/d[123]/e"},
		},
	}
	for _, tt := range tests {
		got := ParseQuery(tt.query)
		if diff := pretty.Compare(tt.parsedQueries, got); diff != "" {
			t.Errorf("ParseQuery(%s) returned diff (-want +got):\n%s", tt.query, diff)
		}
	}
}

func TestParseElement(t *testing.T) {
	tests := []struct {
		element string
		want    []string
	}{
		{
			element: "a",
			want:    []string{"a", ""},
		},
		{
			element: "a[123]",
			want:    []string{"a", "123"},
		},
		{
			element: "a[asd",
			want:    []string{"a[asd", ""},
		},
	}
	for _, tt := range tests {
		gotName, gotKey := parseElement(tt.element)
		got := []string{gotName, gotKey}
		if diff := pretty.Compare(tt.want, got); diff != "" {
			t.Errorf("ToGetRequest(%s) returned diff (-want +got):\n%s", tt.element, diff)
		}
	}
}

func TestToGetRequest(t *testing.T) {
	tests := []struct {
		queries    []string
		getRequest gnmi.GetRequest
	}{
		{
			queries: []string{"a/b/c/d[123]/e", "c/d[123]"},
			getRequest: gnmi.GetRequest{
				Path: []*gnmi.Path{
					&gnmi.Path{
						Elem: []*gnmi.PathElem{
							&gnmi.PathElem{Name: "a"},
							&gnmi.PathElem{Name: "b"},
							&gnmi.PathElem{Name: "c"},
							&gnmi.PathElem{Name: "d", Key: map[string]string{"d": "123"}},
							&gnmi.PathElem{Name: "e"},
						},
					},
					&gnmi.Path{
						Elem: []*gnmi.PathElem{
							&gnmi.PathElem{Name: "c"},
							&gnmi.PathElem{Name: "d", Key: map[string]string{"d": "123"}},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got := ToGetRequest(tt.queries)
		if diff := pretty.Compare(tt.getRequest, got); diff != "" {
			t.Errorf("ToGetRequest(%s) returned diff (-want +got):\n%s", tt.queries, diff)
		}
	}
}
