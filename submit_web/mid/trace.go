package mid

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"submit_web/global"
)

// Trace 生成tracer和startSpan
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
			},
			ServiceName: "submit",
		}

		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}
		opentracing.SetGlobalTracer(tracer)
		//defer closer.Close()

		startSpan := opentracing.StartSpan(c.Request.URL.Path)
		defer startSpan.Finish()

		c.Set("closer", closer)
		c.Set("tracer", tracer)
		c.Set("parentSpan", startSpan)
		c.Next()
	}
}
