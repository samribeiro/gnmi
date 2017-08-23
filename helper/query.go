// Package helper provides helper functions for the gNMI binaries.
package helper

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/santhosh-tekuri/xpathparser"
)

var (
	querySeparator = flag.String("query_separator", ",", "Query separator character.")
)

// ParseQuery parses command line queries.
func ParseQuery(query string) []string {
	var queries []string
	for _, q := range strings.Split(query, *querySeparator) {
		queries = append(queries, q)
	}
	return queries
}

// ToGetRequest generates a gnmi GetRequest out of a list of xPaths.
func ToGetRequest(xpaths []string) (*gnmi.GetRequest, error) {
	getRequest := gnmi.GetRequest{Path: []*gnmi.Path{}}
	for _, xpath := range xpaths {

		expr, err := xpathparser.Parse(xpath)
		if err != nil {
			return nil, err
		}

		locationPath, ok := expr.(*xpathparser.LocationPath)
		if !ok {
			return nil, fmt.Errorf("error parsing LocationPath in xpath")
		}

		path := gnmi.Path{}
		for _, step := range locationPath.Steps {

			nameTest, ok := step.NodeTest.(*xpathparser.NameTest)
			if !ok {
				return nil, fmt.Errorf("error parsing NameTest in xpath")
			}

			pathElem := gnmi.PathElem{Name: nameTest.Local, Key: make(map[string]string)}
			for _, predicate := range step.Predicates {
				binaryExpression, ok := predicate.(*xpathparser.BinaryExpr)
				if !ok {
					return nil, fmt.Errorf("error parsing BinaryExpr in xpath")
				}

				lhs, ok := binaryExpression.LHS.(*xpathparser.LocationPath)
				if !ok {
					return nil, fmt.Errorf("error parsing LHS in xpath")
				}

				if len(lhs.Steps) != 1 {
					return nil, fmt.Errorf("error in LHS length in xpath")
				}

				keyNameTest, ok := lhs.Steps[0].NodeTest.(*xpathparser.NameTest)
				if !ok {
					return nil, fmt.Errorf("error parsing LHS NameTest in xpath")
				}
				key := keyNameTest.Local

				switch rhs := binaryExpression.RHS.(type) {
				case *xpathparser.LocationPath:
					if len(rhs.Steps) != 1 {
						return nil, fmt.Errorf("error in RHS length in xpath")
					}
					valNameTest, ok := rhs.Steps[0].NodeTest.(*xpathparser.NameTest)
					if !ok {
						return nil, fmt.Errorf("error parsing RHS NameTest in xpath")
					}
					pathElem.Key[key] = valNameTest.Local
				case xpathparser.Number:
					pathElem.Key[key] = strconv.FormatFloat(float64(rhs), 'f', -1, 64)
				default:
					return nil, fmt.Errorf("error parsing RHS in xpath")
				}
			}
			path.Elem = append(path.Elem, &pathElem)
		}
		getRequest.Path = append(getRequest.Path, &path)
	}
	return &getRequest, nil
}

// ReflectGetRequest generates a gNMI GetResponse out of a gnmi GetRequest.
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
