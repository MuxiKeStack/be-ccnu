package domain

type Course struct {
	CourseId string `json:"course_id"`
	Name     string `json:"name"`
	Teacher  string `json:"teacher"`
	School   string
	Property string
	Credit   float32
}
