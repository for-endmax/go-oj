package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"go.uber.org/zap"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化sentinel 异常: %v", err)
	}

	//基于错误率
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "submit",
			Strategy:         circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 5,
			StatIntervalMs:   10000,
			MaxAllowedRtMs:   3000,
			Threshold:        0.2,
		},
	})
	if err != nil {
		zap.S().Fatalf("加载规则失败: %v", err)
	}
}
