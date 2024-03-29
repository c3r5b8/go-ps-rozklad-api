package psrozklad

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestGetGroups(t *testing.T) {

	want := []Group{
		{
			Departament: "Фізико-математичний факультет",
			Name:        "11Бд-СОмат",
			Id:          10359,
		},
		{
			Departament: "Фізико-математичний факультет",
			Name:        "11Мд-СОмат",
			Id:          10370,
		},
	}

	api := New("")
	api.HttpClient = &MockHttpClient{resp: http.Response{
		Body: io.NopCloser(bytes.NewBufferString(`{
				"psrozklad_export": {
					"departments": [
						{
							"name": "Фізико-математичний факультет",
							"objects": [
								{
									"name": "11Бд-СОмат",
									"ID": "10359"
								},
								{
									"name": "11Мд-СОмат",
									"ID": "10370"
								}
							]
						}
					],
					"code": "0"
				}
			}`)),
	}}
	got, err := api.GetGroups()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
