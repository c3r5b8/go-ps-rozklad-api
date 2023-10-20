package psrozklad

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestGetRooms(t *testing.T) {
	want := []Room{
		{
			Block:    "гуртож №1",
			Name:     "1",
			FullName: "1/гуртож №1",
			Id:       108,
		},
		{
			Block:    "№1",
			Name:     "320",
			FullName: "320/№1",
			Id:       36,
		},
		{
			Block:    "№1",
			Name:     "Каф. проф. пед.",
			FullName: "Каф. проф. пед./№1",
			Id:       308,
		},
	}
	api := Api{
		HttpClient: &MockHttpClient{
			resp: http.Response{
				Body: io.NopCloser(bytes.NewBufferString(`{
					"psrozklad_export": {
						"blocks": [
							{
								"name": "гуртож №1",
								"objects": [
									{
										"name": "1/гуртож №1",
										"ID": "108"
									}
								]
							},
							{
								"name": "№1",
								"objects": [
									{
										"name": "320/№1",
										"ID": "36"
									},
									{
										"name": "Каф. проф. пед./№1",
										"ID": "308"
									}
								]
							}
						],
						"code": "0"
					}
				}
				`)),
			},
		},
	}
	got, err := api.GetRooms()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
