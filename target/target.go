// Package target provides helper functions for target gNMI binaries.
package target

import "github.com/openconfig/gnmi/proto/gnmi"

// ReflectGetRequest generates a gNMI GetResponse out of a gnmi GetRequest.
func ReflectGetRequest(request *gnmi.GetRequest) *gnmi.GetResponse {
	response := gnmi.GetResponse{Notification: []*gnmi.Notification{}}
	notification := gnmi.Notification{Update: []*gnmi.Update{}}
	for _, path := range request.Path {
		typedValue := gnmi.TypedValue{Value: &gnmi.TypedValue_StringVal{StringVal: "TEST STRING"}}
		update := gnmi.Update{Path: path, Val: &typedValue}
		notification.Update = append(notification.Update, &update)
	}
	response.Notification = append(response.Notification, &notification)
	return &response
}
