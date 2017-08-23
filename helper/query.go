package helper

import (
	"flag"
	"regexp"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
)

var (
	querySeparator = flag.String("query_separator", ",", "Query separator character.")
	queryDelimiter = flag.String("query_delimiter", "/", "Query delimiter character.")
	elementRegex   = regexp.MustCompile(`^([^\[]*)\[([^\]]*)\]$`)
)

func ParseQuery(query string) []string {
	var queries []string
	for _, q := range strings.Split(query, *querySeparator) {
		q := strings.Trim(q, *queryDelimiter)
		queries = append(queries, q)
	}
	return queries
}

func parseElement(element string) (name, key string) {
	subs := elementRegex.FindStringSubmatch(element)
	if len(subs) != 3 {
		return element, ""
	}
	return subs[1], subs[2]
}

func ToGetRequest(queries []string) gnmi.GetRequest {
	getRequest := gnmi.GetRequest{Path: []*gnmi.Path{}}
	for _, query := range queries {
		path := gnmi.Path{}
		for _, element := range strings.Split(query, *queryDelimiter) {
			name, key := parseElement(element)
			pathElem := gnmi.PathElem{}
			if key != "" {
				pathElem.Key = make(map[string]string)
				pathElem.Key[name] = key
			}
			pathElem.Name = name
			path.Elem = append(path.Elem, &pathElem)

		}
		getRequest.Path = append(getRequest.Path, &path)
	}
	return getRequest
}

func ReflectGetRequest(request *gnmi.GetRequest) *gnmi.GetResponse {
	response := gnmi.GetResponse{Notification: []*gnmi.Notification{}}
	notification := gnmi.Notification{Update: []*gnmi.Update{}}
	for _, path := range request.Path {
		typedValue := gnmi.TypedValue{Value: &gnmi.TypedValue_StringVal{StringVal: "TESTDATA"}}
		update := gnmi.Update{Path: path, Val: &typedValue}
		notification.Update = append(notification.Update, &update)
	}
	response.Notification = append(response.Notification, &notification)
	return &response
}
