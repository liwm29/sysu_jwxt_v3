package course

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/liwm29/sysu_jwxt_v3/backend/jwxtClient/request"
	"time"
)

// 如果作为一个三方库设计,不应该是一个完全黑盒,所以这里也传递选课失败信息
func (c *Course) AutoChoose(cl request.Clienter, sleep time.Duration) <-chan error {
	isOk := make(chan error)
	go func() {
		color.Green("开始选课任务,刷新间隔:%#v seconds", sleep/time.Second)
		for {
			check(c.Refresh(cl))

			if c.IsSelected() {
				isOk <- nil
				break
			}

			if c.VacancyNum() == 0 {
				isOk <- errors.New("选课失败: " + c.VacancyInfo())
				time.Sleep(sleep)
				continue
			}

			if ok := c.Choose(cl); ok {
				check(c.Refresh(cl))
				if c.IsSelected() {
					isOk <- nil
					break
				}
			}
		}
		close(isOk)
	}()
	return isOk
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
