package service

import (
	"context"
	"encoding/json"
	"github.com/MuxiKeStack/be-ccnu/domain"
	"github.com/ecodeclub/ekit/slice"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GradeList struct {
	Items []GradeItem `json:"items"`
}

type GradeItem struct {
	//JgID               string `json:"jg_id"`
	Jsxm   string `json:"jsxm"`   // 教工名称
	Kch    string `json:"kch"`    // 课程号
	Kcmc   string `json:"kcmc"`   // 课程名称
	Kcxzmc string `json:"kcxzmc"` // 课程性质名称
	Kkbmmc string `json:"kkbmmc"` // 开课学院
	Xf     string `json:"xf"`     // 学分
	Cj     string `json:"cj"`     // 成绩
	JxbId  string `json:"jxb_id"`
	Xnm    string `json:"xnm"` // 学年名
	Xqm    string `json:"xqm"` // 学期名
}

func (c *ccnuService) GetSelfGradeList(ctx context.Context, studentId, password, year, term string) ([]domain.Grade, error) {
	client, err := c.xkLoginClient(ctx, studentId, password) // 登录，直接
	if err != nil {
		return nil, err
	}
	var termMap = map[string]string{"1": "3", "2": "12", "3": "16"} // 学期参数
	if year == "0" {
		year = ""
	}
	formData := url.Values{}
	formData.Set("xnm", year)          // 学年名
	formData.Set("xqm", termMap[term]) // 学期名
	formData.Set("_search", "false")
	formData.Set("nd", string(time.Now().UnixNano()))
	formData.Set("queryModel.showCount", "1000")
	formData.Set("queryModel.currentPage", "1")
	formData.Set("queryModel.sortName", "")
	formData.Set("queryModel.sortOrder", "asc")
	formData.Set("time", "5")

	requestUrl := "http://xk.ccnu.edu.cn/jwglxt/cjcx/cjcx_cxXsgrcj.html?doType=query&gnmkdm=N305005"
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0")
	req.Header.Set("Referer", "http://xk.ccnu.edu.cn/jwglxt/cjcx/cjcx_cxDgXscj.html?gnmkdm=N305005&layout=default")
	req.Header.Set("Origin", "http://xk.ccnu.edu.cn")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var gl GradeList
	err = json.Unmarshal(body, &gl)
	if err != nil {
		return nil, err
	}

	res := slice.Map(gl.Items, func(idx int, src GradeItem) domain.Grade {
		credit, _ := strconv.ParseFloat(src.Xf, 10)
		total, _ := strconv.ParseFloat(src.Cj, 10)
		return domain.Grade{
			Course: domain.Course{
				CourseId: src.Kch,
				Name:     src.Kcmc,
				Teacher:  src.Jsxm,
				School:   src.Kkbmmc,
				Property: src.Kcxzmc,
				Credit:   float32(credit),
			},
			Total: float32(total),
			Year:  year,
			Term:  term,
			JxbId: src.JxbId,
			Xnm:   src.Xnm,
			Xqm:   src.Xqm,
		}
	})
	return res, nil
}