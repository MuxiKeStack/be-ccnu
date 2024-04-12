package domain

type Course struct {
	CourseId string `json:"course_id"`
	Name     string `json:"name"`
	Teacher  string `json:"teacher"`
	School   string // 开课学院
	Property string
	Credit   float32
	Year     string
	Term     string
}

type Grade struct {
	Course  Course
	Regular float32
	Final   float32
	Total   float32
	// 下面三个字段和Course中的School可用于进一步查询平时分的期末分
	Year  string
	Term  string
	JxbId string `json:"jxb_id"`
}
