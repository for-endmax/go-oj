package handler

// 自定义 amqp.Table 实现 opentracing.TextMapCarrier 接口
type amqpTableCarrier struct {
	headers map[string]interface{}
}

func (c amqpTableCarrier) Set(key, val string) {
	c.headers[key] = val
}

func (c amqpTableCarrier) ForeachKey(handler func(key, val string) error) error {
	for key, val := range c.headers {
		valStr, ok := val.(string)
		if !ok {
			continue
		}
		if err := handler(key, valStr); err != nil {
			return err
		}
	}
	return nil
}
