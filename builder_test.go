package httplayer

import (
	"net/http"
	"reflect"
	"testing"
)

func TestRoute_GetMethod(t *testing.T) {
	type fields struct {
		method      string
		path        string
		handlerFunc http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "t1", fields: fields{method: "GET", path: "/api/test", handlerFunc: nil}, want: []string{"GET"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				method:      []string{tt.fields.method},
				path:        tt.fields.path,
				handlerFunc: tt.fields.handlerFunc,
			}
			if got := r.Methods(); reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.Method() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_GetPath(t *testing.T) {
	type fields struct {
		method      string
		path        string
		handlerFunc http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "t1", fields: fields{method: "GET", path: "/api/test", handlerFunc: nil}, want: "/api/test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				method:      []string{tt.fields.method},
				path:        tt.fields.path,
				handlerFunc: tt.fields.handlerFunc,
			}
			if got := r.Path(); got != tt.want {
				t.Errorf("Route.Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Build(t *testing.T) {
	r := NewBuilder("GET").Handler(simpleHandler).Build()
	rr, err := testResponse(r.handlerFunc)
	if err != nil {
		t.Fatalf("got error on testresponse: %v", err)
	}
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
