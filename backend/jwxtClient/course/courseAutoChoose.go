package course

import (
	"errors"
	"server/backend/jwxtClient/request"
	"server/backend/jwxtClient/util"
	"strings"
	"time"
)

// 如果作为一个三方库设计,不应该是一个完全黑盒,所以这里也传递选课失败信息
func (c *Course) AutoChoose(cl request.Clienter, sleep time.Duration) <-chan error {
	req := NewCourseListReq(c.CourseType, c.Option).SetCourseName(c.BaseInfo.CourseNum)
	isOk := make(chan error)
	go func() {
		for {
			courseList, n_page := req.DoPage(cl)
			if n_page != 1 || len(courseList.Courses) != 1 {
				log.WithField("n_page", n_page).Info("注意:查询到不止一个页面或不止一个课程\t", util.WhereAmI())
			}

			course, err := courseList.First()
			if err != nil {
				isOk <- err
				break
			}
			if course.courseId() != c.courseId() {
				isOk <- errors.New("课程不存在: " + strings.Join(courseList.CourseNames(), ","))
			}

			if course.VacancyNum() == 0 {
				isOk <- errors.New("选课失败: " + course.VacancyInfo())
				time.Sleep(sleep)
				continue
			}

			if ok := course.Choose(cl); ok {
				isOk <- nil
				break
			}
		}

		close(isOk)
	}()
	return isOk
}
