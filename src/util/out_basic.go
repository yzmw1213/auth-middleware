package util

type OutputBasic struct {
	Code    int         // コード
	Result  string      // 結果
	Message interface{} // メッセージ
}

func (o *OutputBasic) GetResult() map[string]interface{} {
	switch value := o.Message.(type) {
	case error:
		return map[string]interface{}{
			"code":    o.Code,
			"result":  o.Result,
			"message": value.Error(),
		}
	}
	return map[string]interface{}{
		"code":    o.Code,
		"result":  o.Result,
		"message": o.Message,
	}
}

func (o *OutputBasic) GetCode() int {
	return o.Code
}

func (o *OutputBasic) GetError() error {
	switch value := o.Message.(type) {
	case error:
		return value
	}
	return nil
}
