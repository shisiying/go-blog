package global

import (
	"github.com/shisiying/go-blog/pkg/logger"
	"github.com/shisiying/go-blog/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS

	Logger       *logger.Logger
	JWTSetting   *setting.JWTSetting
	EmailSetting *setting.EmailSettingS
)
