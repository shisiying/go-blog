package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shisiying/go-blog/global"
	"github.com/shisiying/go-blog/pkg/app"
	"github.com/shisiying/go-blog/pkg/email"
	"github.com/shisiying/go-blog/pkg/errcode"
	"time"
)

func Recovery() gin.HandlerFunc {

	defailtMailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		From:     global.EmailSetting.From,
	})
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.WithCallersFrames().Errorf(c, "panic recover err: %v", err)
				app.NewResponse(c).ToErrorResponse(errcode.ServerError)

				err := defailtMailer.SendMail(
					global.EmailSetting.To,
					fmt.Sprintf("异常跑出，发生时间：%d", time.Now().Unix()),
					fmt.Sprintf("错误信息：%v", err),
				)

				if err != nil {
					global.Logger.Panicf(c, "mail.SendMail err:%v", err)
				}

				app.NewResponse(c).ToErrorResponse(errcode.ServerError)

				c.Abort()
			}
		}()
		c.Next()
	}
}
