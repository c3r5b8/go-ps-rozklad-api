package psrozklad

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetLessons(t *testing.T) {
	testCases := []struct {
		desc        string
		json        string
		teacherJson string
		groupeJson  string
		roomJson    string

		want []Lesson
	}{
		{
			desc: "test simple lesson",
			json: `{
				"psrozklad_export": {
					"roz_items": [
						{
							"object": "26Бд-Комп",
							"date": "16.10.2023",
							"comment": "0",
							"lesson_number": "1",
							"lesson_name": "1",
							"lesson_time": "09:00-10:20",
							"half": "",
							"teacher": "Горобець С.М.",
							"teachers_add": "",
							"room": "320/№1",
							"group": "Потік 21Бд-СОмат, 22Бд-СОмат",
							"title": "Інженерна та комп‘ютерна графіка",
							"type": "Л",
							"replacement": "",
							"reservation": "",
							"online": "Так",
							"comment4link": "Ідентифікатор: 979 971 2364; Пароль: 2023",
							"link": "https://us05web.zoom.us/j/9799712364?pwd=f5hSQnbCbnvU6ACFWEyQT6wMBBzk0v.1"
						}
					]
				}
			}`,
			teacherJson: `{
				"psrozklad_export": {
					"departments": [
						{
							"name": "Кафедра комп‘ютерних наук та інформаційних технологій",
							"objects": [
								{
									"name": "Горобець С.М.",
									"P": "Горобець",
									"I": "C",
									"B": "Миколайович",
									"ID": "420"
								}
							]
						}
					],
					"code": "0"
				}
			}`,
			groupeJson: `{
				"psrozklad_export": {
					"departments": [
						{
							"name": "Фізико-математичний факультет",
							"objects": [
								{
									"name": "21Бд-СОмат",
									"ID": "11"
								},
								{
									"name": "22Бд-СОмат",
									"ID": "12"
								}
							]
						}
					],
					"code": "0"
				}
			}`,
			roomJson: `{
				"psrozklad_export": {
					"blocks": [
						{
							"name": "№1",
							"objects": [
								{
									"name": "320/№1",
									"ID": "36"
								}
							]
						}
					],
					"code": "0"
				}
			}`,
			want: []Lesson{
				{
					Title: "Інженерна та комп‘ютерна графіка",
					Teacher: Teacher{
						ShortName:   "Горобець С.М.",
						P:           "Горобець",
						I:           "C",
						B:           "Миколайович",
						Id:          420,
						Departament: "Кафедра комп‘ютерних наук та інформаційних технологій",
					},
					Type:   "Л",
					Day:    "16.10.2023",
					Number: 1,
					Room: Room{
						Block:    "№1",
						Name:     "320",
						FullName: "320/№1",
						Id:       36,
					},
					GroupsType: "Потік",
					Groups: []Group{
						{
							Name:        "21Бд-СОмат",
							Id:          11,
							Departament: "Фізико-математичний факультет",
						},
						{
							Name:        "22Бд-СОмат",
							Id:          12,
							Departament: "Фізико-математичний факультет",
						},
					},
					SubGroup:       "21Бд-СОмат, 22Бд-СОмат",
					StartTime:      time.Date(2023, time.October, 16, 9, 0, 0, 0, time.Local),
					EndTime:        time.Date(2023, time.October, 16, 10, 20, 0, 0, time.Local),
					Online:         true,
					URL:            "https://us05web.zoom.us/j/9799712364?pwd=f5hSQnbCbnvU6ACFWEyQT6wMBBzk0v.1",
					CommentForLink: "Ідентифікатор: 979 971 2364; Пароль: 2023",
				},
			},
		},
		{
			desc: "test sub groups",
			json: `{
				"psrozklad_export": {
					"roz_items": [
						{
							"object": "22Бд-СОмат",
							"date": "16.10.2023",
							"comment": "0",
							"lesson_number": "4",
							"lesson_name": "4",
							"lesson_time": "13:40-15:00",
							"half": "",
							"teacher": "Яценко О.С.",
							"teachers_add": "",
							"room": "320/№1",
							"group": "(підгр. 1)",
							"title": "Комп‘ютерні мережі",
							"type": "Лаб",
							"replacement": "",
							"reservation": "",
							"online": "",
							"comment4link": "",
							"link": ""
						},
						{
							"object": "22Бд-СОмат",
							"date": "16.10.2023",
							"comment": "0",
							"lesson_number": "4",
							"lesson_name": "4",
							"lesson_time": "13:40-15:00",
							"half": "",
							"teacher": "Кривонос О.М.",
							"teachers_add": "",
							"room": "320/№1",
							"group": "Збірна група 21Бд-СОмат, 22Бд-СОмат",
							"title": "Комп‘ютерні мережі",
							"type": "Лаб",
							"replacement": "",
							"reservation": "",
							"online": "",
							"comment4link": "",
							"link": ""
						}
					],
					"code": "0"
				}
			}
			`,
			teacherJson: `{
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
			}`,
			roomJson: `{
				"psrozklad_export": {
					"blocks": [
						{
							"name": "№1",
							"objects": [
								{
									"name": "320/№1",
									"ID": "36"
								}
							]
						}
					],
					"code": "0"
				}
			}`,
			groupeJson: `{
				"psrozklad_export": {
					"departments": [
						{
							"name": "Фізико-математичний факультет",
							"objects": [
								{
									"name": "21Бд-СОмат",
									"ID": "11"
								},
								{
									"name": "22Бд-СОмат",
									"ID": "12"
								}
							]
						}
					],
					"code": "0"
				}
			}`,
			want: []Lesson{
				{
					Title: "Комп‘ютерні мережі",
					Teacher: Teacher{
						ShortName:   "Яценко О.С.",
						P:           "Яценко",
						I:           "Олександр",
						B:           "Сергійович",
						Departament: "Кафедра комп‘ютерних наук та інформаційних технологій",
						Id:          486,
					},
					Type:   "Лаб",
					Day:    "16.10.2023",
					Number: 4,
					Room: Room{
						Block:    "№1",
						Name:     "320",
						FullName: "320/№1",
						Id:       36,
					},
					GroupsType: "підгр",
					Groups: []Group{
						{
							Name:        "22Бд-СОмат",
							Id:          12,
							Departament: "Фізико-математичний факультет",
						},
					},
					SubGroup:  "(підгр. 1)",
					StartTime: time.Date(2023, time.October, 16, 13, 40, 0, 0, time.Local),
					EndTime:   time.Date(2023, time.October, 16, 15, 0, 0, 0, time.Local),
				},
				{
					Title: "Комп‘ютерні мережі",
					Teacher: Teacher{
						ShortName:   "Кривонос О.М.",
						P:           "Кривонос",
						I:           "Олександр",
						B:           "Миколайович",
						Departament: "Кафедра комп‘ютерних наук та інформаційних технологій",
						Id:          420,
					},
					Type:   "Лаб",
					Day:    "16.10.2023",
					Number: 4,
					Room: Room{
						Block:    "№1",
						Name:     "320",
						FullName: "320/№1",
						Id:       36,
					},
					GroupsType: "підгр",
					Groups: []Group{
						{
							Name:        "21Бд-СОмат",
							Id:          11,
							Departament: "Фізико-математичний факультет",
						},
						{
							Name:        "22Бд-СОмат",
							Id:          12,
							Departament: "Фізико-математичний факультет",
						},
					},
					SubGroup:  "21Бд-СОмат, 22Бд-СОмат",
					StartTime: time.Date(2023, time.October, 16, 13, 40, 0, 0, time.Local),
					EndTime:   time.Date(2023, time.October, 16, 15, 0, 0, 0, time.Local),
				},
			},
		},
		// {
		// 	desc: "test_teacher",
		// 	teacherJson: `{
		// 		"psrozklad_export": {
		// 			"departments": [
		// 				{
		// 					"name": "Кафедра комп‘ютерних наук та інформаційних технологій",
		// 					"objects": [
		// 						{
		// 							"name": "Кривонос О.М.",
		// 							"P": "Кривонос",
		// 							"I": "Олександр",
		// 							"B": "Миколайович",
		// 							"ID": "420"
		// 						},
		// 						{
		// 							"name": "Яценко О.С.",
		// 							"P": "Яценко",
		// 							"I": "Олександр",
		// 							"B": "Сергійович",
		// 							"ID": "486"
		// 						}
		// 					]
		// 				}
		// 			],
		// 			"code": "0"
		// 		}
		// 	}`,
		// 	// roomJson: ,
		// },
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			api := Api{
				HttpClient: &MockHttpClient{
					http.Response{
						Body: io.NopCloser(bytes.NewBufferString(tC.json)),
					},
				},
			}
			apiRooms := Api{
				HttpClient: &MockHttpClient{
					http.Response{
						Body: io.NopCloser(bytes.NewBufferString(tC.roomJson)),
					},
				},
			}
			apiRooms.InitRooms()
			apiTeachers := Api{
				HttpClient: &MockHttpClient{
					http.Response{
						Body: io.NopCloser(bytes.NewBufferString(tC.teacherJson)),
					},
				},
			}
			apiTeachers.InitTeachers()
			apiGroups := Api{
				HttpClient: &MockHttpClient{
					http.Response{
						Body: io.NopCloser(bytes.NewBufferString(tC.groupeJson)),
					},
				},
			}
			apiGroups.InitGroups()

			got, err := api.GetLessons(Group{}, time.Time{}, time.Time{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Errorf("want: \n%v\ngot: \n%v", tC.want, got)
			}

		})
	}
}
