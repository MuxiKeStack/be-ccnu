package domain

type Course struct {
	CourseId string `json:"course_id"`
	Name     string `json:"name"`
	Teacher  string `json:"teacher"`
	Year     string `json:"year"` // 学期，2018
	Term     string `json:"term"` // 学年，1/2/3
}
