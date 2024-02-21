package httplayer

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestRoutingDefinition(t *testing.T) {
	type fields struct {
		values []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:    "without parent middleware",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "one parent middleware",
			fields: fields{
				values: []string{"0"},
			},
			want: []string{"0"},
		},
		{
			name: "many parent middleware",
			fields: fields{
				values: []string{"0", "1", "4", "8"},
			},
			want: []string{"0", "1", "4", "8"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestRoute(tt.fields.values...).handlerFunc
			rr, err := testResponse(h)
			if err != nil {
				t.Fatal(err)
			}
			isOk := rr.Code == http.StatusOK
			if !isOk {
				if !tt.wantErr {
					t.Fatalf("handler returned wrong status code: got %d want %d",
						rr.Code, http.StatusOK)
				}
				return
			}
			var calls []string
			if err := json.NewDecoder(rr.Body).Decode(&calls); err != nil {
				t.Fatalf("could not decode response")
			}
			if !reflect.DeepEqual(calls, tt.want) {
				t.Fatalf("expected response %v, got %v", tt.want, calls)
			}
		})
	}
}

func TestRoutingDefinitionWithChildMiddlewares(t *testing.T) {
	type fields struct {
		parentValues []string
		childValues  []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:    "without middleware",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "only child middleware",
			fields: fields{
				childValues: []string{"child"},
			},
			want: []string{"child"},
		},
		{
			name: "one parent middleware - one child",
			fields: fields{
				parentValues: []string{"1"},
				childValues:  []string{"0"},
			},
			want: []string{"1", "0"},
		},
		{
			name: "many parent middleware - many child",
			fields: fields{
				parentValues: []string{"0", "1", "4", "8"},
				childValues:  []string{"10", "11"},
			},
			want: []string{"0", "1", "4", "8", "10", "11"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestRouteWithMws(tt.fields.parentValues, tt.fields.childValues).handlerFunc
			rr, err := testResponse(h)
			if err != nil {
				t.Fatal(err)
			}
			isOk := rr.Code == http.StatusOK
			if !isOk {
				if !tt.wantErr {
					t.Fatalf("handler returned wrong status code: got %d want %d",
						rr.Code, http.StatusOK)
				}
				return
			}
			var calls []string
			if err := json.NewDecoder(rr.Body).Decode(&calls); err != nil {
				t.Fatalf("could not decode response")
			}
			if !reflect.DeepEqual(calls, tt.want) {
				t.Fatalf("expected response %v, got %v", tt.want, calls)
			}
		})
	}
}
