package psrozklad

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Object interface {
	ID() int
	type_obj() string
}

type Lesson struct {
	Title          string
	Teacher        Teacher
	Type           string
	Day            string
	Number         int
	Room           Room
	GroupsType     string
	Groups         []Group
	SubGroup       string
	StartTime      time.Time
	EndTime        time.Time
	Online         bool
	URL            string
	CommentForLink string
}

type lessonExport struct {
	Object       string `json:"object"`
	Date         string `json:"date"`
	Lesson_time  string `json:"lesson_time"`
	Teacher      string `json:"teacher"`
	Number       string `json:"lesson_number"`
	Room         string `json:"room"`
	Group        string `json:"group"`
	Title        string `json:"title"`
	Replacement  string `json:"replacement"`
	Type         string `json:"type"`
	Online       string `json:"online"`
	Link         string `json:"link"`
	Comment4link string `json:"comment4link"`
}

// GetLessons gets the lessons from the timetable export for the given object and time period.
func (a *Api) GetLessons(obj Object, start, end time.Time) ([]Lesson, error) {

	// Create a timetable export struct to decode the JSON response into.
	type timetableExport struct {
		RozItems []lessonExport `json:"roz_items"`
		Code     string         `json:"code"`
	}

	// Create an export struct to decode the JSON response into.
	type export struct {
		Timetable timetableExport `json:"psrozklad_export"`
	}

	// Build the URL for the timetable export request.
	url := a.BaseUri
	startY, startM, startD := start.Date()
	endY, endM, endD := end.Date()
	url += "&begin_date=" + strconv.Itoa(startD) + "." + strconv.Itoa(int(startM)) + "." + strconv.Itoa(startY)
	url += "&end_date=" + strconv.Itoa(endD) + "." + strconv.Itoa(int(endM)) + "." + strconv.Itoa(endY)
	url += "&OBJ_ID=" + strconv.Itoa(obj.ID()) + "&ros_text=separated"
	url += "&req_mode=" + obj.type_obj()

	// Create a new HTTP request for the timetable export.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Make the request to the timetable export endpoint.
	resp, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get lessons: %v", err)
	}
	defer resp.Body.Close()

	// Decode the JSON response into an export struct.
	var exp export
	err = json.NewDecoder(resp.Body).Decode(&exp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode lessons json: %v", err)
	}

	// Create a slice to store the lessons.
	var lessons, lessons_temp []Lesson

	// Keep track of the current lesson number.
	var num int

	// Iterate over the lesson export items in the timetable export.
	for _, lesson := range exp.Timetable.RozItems {

		// Convert the lesson export item to a Lesson struct.
		less_new, err := convertLessonExportToLesson(lesson, obj.type_obj())
		if err != nil {
			return nil, fmt.Errorf("failed to convert lessonExport to Lesson: %v", err)
		}

		// If the current lesson number is 0, add the lesson to the lessons_temp slice.
		if num == 0 {
			lessons_temp = append(lessons_temp, less_new)
		} else if less_new.Number == num {
			// If the current lesson number is the same as the previous lesson number, add the lesson to the lessons_temp slice.
			lessons_temp = append(lessons_temp, less_new)
		} else if less_new.Number != num {
			// If the current lesson number is different from the previous lesson number, convert the lessons_temp slice to a slice of Lesson structs and add it to the lessons slice. Then, clear the lessons_temp slice and add the current lesson to it.
			lessons = append(lessons, convertLessons(lessons_temp)...)
			lessons_temp = nil
			lessons_temp = append(lessons_temp, less_new)
		}
	}

	// Convert the lessons_temp slice to a slice of Lesson structs and add it to the lessons slice.
	lessons = append(lessons, convertLessons(lessons_temp)...)

	// Return the lessons slice.
	return lessons, nil
}
