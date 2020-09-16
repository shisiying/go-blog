package routers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/shisiying/go-blog/docs"
	"github.com/shisiying/go-blog/global"
	"github.com/shisiying/go-blog/internal/middleware"
	"github.com/shisiying/go-blog/internal/routers/api"
	"github.com/shisiying/go-blog/internal/routers/api/v1"
	"github.com/shisiying/go-blog/pkg/limiter"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"time"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	r := gin.New()

	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTimeout(60 * time.Second))

	r.Use(middleware.Translations())
	r.Use(middleware.JWT())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/auth", api.GetAuth)

	tag := v1.NewTag()
	article := v1.NewArticle()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles/", article.List)
	}

	return r
}
