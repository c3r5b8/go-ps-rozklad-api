package psrozklad

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// convertLessonExportToLesson converts a lessonExport struct to a Lesson struct and returns it along with an error.
// It takes the lessonExport and a type as input parameters.
func convertLessonExportToLesson(les lessonExport, t string) (Lesson, error) {
	// Initialize the error variable.
	var err error

	// Depending on the value of type, update the relevant fields of the `les` struct.
	switch t {
	case "room":
		// If `t` is "room," set the Room field in `les` to the Object field.
		les.Room = les.Object
	case "teacher":
		// If `t` is "teacher," split the Object into parts, format it as a teacher's name, and assign it to the Teacher field in `les`.
		teacher_list := strings.Split(les.Object, " ")
		les.Teacher = teacher_list[0] + " " + string(teacher_list[1][:2]) + "." + string(teacher_list[2][:2]) + "."
	}

	// Create a new Lesson struct and initialize it with some fields from the `les` struct.
	less_new := Lesson{
		Title:   les.Title,
		Teacher: Teachers[les.Teacher],
		Type:    les.Type,
	}
	less_new.Room = Rooms[les.Room]

	// Use the convertTime function to parse lesson time and date, and assign the result to the StartTime and EndTime fields in `less_new`.
	less_new.StartTime, less_new.EndTime, err = convertTime(les.Lesson_time, les.Date)
	if err != nil {
		return Lesson{}, fmt.Errorf("failed to convert lesson time: %v", err)
	}

	// If `t` is not "group" and the Group field in `les` ends with a closing parenthesis, split the Group field to populate the Groups, GroupsType, and SubGroup fields.
	if t != "group" && strings.HasSuffix(les.Group, ")") {
		splited_groupe := strings.Split(les.Group, " ")
		less_new.Groups, less_new.GroupsType, less_new.SubGroup = convertStringToGroupes(splited_groupe[0], splited_groupe[1]+" "+splited_groupe[2])
	} else {
		// Otherwise, use the Object and Group fields to populate the Groups, GroupsType, and SubGroup fields.
		less_new.Groups, less_new.GroupsType, less_new.SubGroup = convertStringToGroupes(les.Object, les.Group)
	}

	// Convert the Number field in `les` to an integer and assign it to the Number field in `less_new`.
	less_new.Number, err = strconv.Atoi(les.Number)
	if err != nil {
		return Lesson{}, err
	}

	// Set the Day field in `less_new` to the Date field in `les`.
	less_new.Day = les.Date

	// Check if the Online field in `les` is "Так" (meaning "Yes"), and if so, set relevant fields in `less_new`.
	if les.Online == "Так" {
		less_new.Online = true
		less_new.URL = les.Link
		less_new.CommentForLink = les.Comment4link
	}

	// If there is a Replacement field in `les`, handle it by Replacement field in `less_new`.
	if les.Replacement != "" {
		less_new.Replacement.Teacher, less_new.Replacement.Title, less_new.Replacement.Type = convertReplacment(les.Replacement)
	}

	// Return the converted Lesson struct and a nil error.
	return less_new, nil
}

// convertReplacment converts a replacement string to a Teacher struct, a title string, and a lesson type string.
func convertReplacment(replacment string) (Teacher Teacher, title string, leson_typ string) {
	// If the replacement string ends with "замість:", then it is a replacement for a teacher and a lesson.
	if strings.HasSuffix(replacment, "замість:") {
		// Get rid of the "Увага! Заміна! " and " замість:" prefixes.
		replacment = replacment[26 : len(replacment)-16]

		// Split the replacement string into a slice of strings.
		replacmentSlice := strings.Split(replacment, " ")

		// The first two strings in the slice are the teacher's name.
		teacher := replacmentSlice[0]
		teacher = strings.ReplaceAll(teacher, " ", " ")
		Teacher = Teachers[teacher]

		// The last string in the slice is the lesson type.
		leson_typ = replacmentSlice[len(replacmentSlice)-1]

		// Concatenate the remaining strings in the slice to form the title.
		for i := 1; i < len(replacmentSlice)-2; i++ {
			title += replacmentSlice[i] + " "
		}
		title += replacmentSlice[len(replacmentSlice)-2]
	} else {
		// If the replacement string does not end with "замість:", then it is a replacement for a teacher only.
		replacmentSlice := strings.Split(replacment, ": ")

		teacher := replacmentSlice[0][26 : len(replacmentSlice[0])-15]
		teacher = strings.ReplaceAll(teacher, " ", " ")
		Teacher = Teachers[teacher]
	}

	// Return the Teacher struct, the title string, and the lesson type string.
	return Teacher, title, leson_typ
}

// Convert a string representing a date and time to two time.Time objects.
func convertTime(t, d string) (time.Time, time.Time, error) {
	// Split the date string into three parts: year, month, and day.
	dataa := strings.Split(d, ".")

	// Split the time string into two parts: start time and end time.
	times := strings.Split(t, "-")

	// Create a 2D slice to store the start and end times.
	var tt [][]string
	for _, v := range times {
		tt = append(tt, strings.Split(v, ":"))
	}

	// Convert the year, month, day, start hour, start minute, end hour, and end minute to integers.
	year, err := strconv.Atoi(dataa[2])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	month, err := strconv.Atoi(dataa[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	day, err := strconv.Atoi(dataa[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start_h, err := strconv.Atoi(tt[0][0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start_m, err := strconv.Atoi(tt[0][1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end_h, err := strconv.Atoi(tt[1][0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end_m, err := strconv.Atoi(tt[1][1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Create two time.Time objects for the start and end times.
	start := time.Date(year, time.Month(month), day, start_h, start_m, 0, 0, time.Local)
	end := time.Date(year, time.Month(month), day, end_h, end_m, 0, 0, time.Local)

	// Return the start and end time.Time objects.
	return start, end, nil
}

// convertStringToGroupes converts a string representation of a group to a slice of Group
// objects and a string representing the group type.
func convertStringToGroupes(obj, groups string) ([]Group, string, string) {
	var gr_l []Group
	var gr_t, subGr string

	// If the groups string is empty, then the group is the same as the obj string.
	if groups == "" {
		gr_l = append(gr_l, Groups[obj])
	} else if strings.HasPrefix(groups, "(") {
		// If the groups string starts with a parenthesis, then it is a subgroup.
		gr_l = append(gr_l, Groups[obj])
		subGr = groups
		gr_t = "підгр"
	} else if strings.HasPrefix(groups, "Збірна група") {
		// If the groups string starts with "Збірна група", then it is a combined group.
		s := groups[24:]
		subGr = s
		groups := strings.Split(s, ", ")
		for _, gr := range groups {
			gr_l = append(gr_l, Groups[gr])
		}
		gr_t = "Збірна група"
	} else if strings.HasPrefix(groups, "Потік") {
		// If the groups string starts with "Потік", then it is a stream.
		s := groups[11:]
		subGr = s
		groups := strings.Split(s, ", ")
		for _, gr := range groups {
			gr_l = append(gr_l, Groups[gr])
		}
		gr_t = "Потік"
	} else {
		// Otherwise, the group is the same as the groups string.
		gr_l = append(gr_l, Groups[groups])
	}

	return gr_l, gr_t, subGr
}

// convertLessons converts all lessons with a GroupsType of "підгр" to have a GroupsType of "підгр".
func convertLessons(lessons []Lesson) []Lesson {
	// ifSubGroupe is a flag that indicates whether any of the lessons have a GroupsType of "підгр".
	var ifSubGroupe bool
	for i := 0; i < len(lessons); i++ {
		if lessons[i].GroupsType == "підгр" {
			ifSubGroupe = true
		}
	}

	// If any of the lessons have a GroupsType of "підгр", then convert all lessons to have a GroupsType of "підгр".
	if ifSubGroupe {
		for i := 0; i < len(lessons); i++ {
			lessons[i].GroupsType = "підгр"
		}
	}

	return lessons
}
