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

	//基于慢请求
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "submit",
			Strategy:         circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:   10000,
			MinRequestAmount: 2, //静默
			StatIntervalMs:   20000,
			MaxAllowedRtMs:   2000,
			Threshold:        0.4,
		},
	})
	if err != nil {
		zap.S().Fatalf("加载规则失败: %v", err)
	}
}
