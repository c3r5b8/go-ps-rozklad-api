package psrozklad

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestGetTeachers(t *testing.T) {
	want := []Teacher{
		{
			ShortName:   "Кривонос О.М.",
			P:           "Кривонос",
			I:           "Олександр",
			B:           "Миколайович",
			Departament: "Кафедра комп‘ютерних наук та інформаційних технологій",
			Id:          420,
		},
		{
			ShortName:   "Яценко О.С.",
			P:           "Яценко",
			I:           "Олександр",
			B:           "Сергійович",
			Departament: "Кафедра комп‘ютерних наук та інформаційних технологій",
			Id:          486,
		},
	}
	api := Api{
		HttpClient: &MockHttpClient{
			resp: http.Response{
				Body: io.NopCloser(bytes.NewBufferString(`{
					"psrozklad_export": {
						"departments": [
							{
								"name": "Кафедра комп‘ютерних наук та інформаційних технологій",
								"objects": [
									{
										"name": "Кривонос О.М.",
										"P": "Кривонос",
										"I": "Олександр",
										"B": "Миколайович",
										"ID": "420"
									},
									{
										"name": "Яценко О.С.",
										"P": "Яценко",
										"I": "Олександр",
										"B": "Сергійович",
										"ID": "486"
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
	got, err := api.GetTeachers()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
