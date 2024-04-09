package domain

type Course struct {
	CourseId string `json:"course_id"`
	Name     string `json:"name"`
	Teacher  string `json:"teacher"`
	School   string // 开课学院
	Property string
	Credit   float32
}

type Grade struct {
	Course  Course
	Regular float32
	Final   float32
	Total   float32
	Year    string
	Term    string
	// 下面三个字段和Course中的School用于进一步查询平时分的期末分
	JxbId string `json:"jxb_id"`
	Xnm   string `json:"xnm"` // 学年名
	Xqm   string `json:"xqm"` // 学期名
}
