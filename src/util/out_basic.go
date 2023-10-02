package util

import (
	"net/http"
	"reflect"
)

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

type OutputBasicInterface interface {
	GetCode() int
	GetResult() map[string]interface{}
}

type OutputBasicServiceInterface interface {
	GetCode() int
	GetResult() map[string]interface{}
	GetError() error
}

type OutputBasicObject struct {
	Code    int          // コード
	Result  string       // 結果
	OutData *interface{} // 結果オブジェクト
}

func NewOutputBasicObject(row interface{}) *OutputBasicObject {
	return &OutputBasicObject{
		OutData: &row,
	}
}

func (o *OutputBasicObject) GetResult() map[string]interface{} {
	value := reflect.ValueOf(*o.OutData)
	if value.Kind() == reflect.Slice {
		return map[string]interface{}{
			"code":   http.StatusOK,
			"result": "OK",
			"list":   o.OutData,
		}
	} else {
		return map[string]interface{}{
			"code":   http.StatusOK,
			"result": "OK",
			"row":    o.OutData,
		}
	}
}

func (o *OutputBasicObject) GetCode() int {
	return o.Code
}

type OutputBasicListPaging struct {
	Code       int          // コード
	Result     string       // 結果
	List       *interface{} // 結果オブジェクト配列
	CountTotal int64        // 総数
	Page       int64        // ページ番号
	Limit      int64        // 検索数
}

func NewOutputBasicListPaging(row interface{}) *OutputBasicListPaging {
	return &OutputBasicListPaging{
		List: &row,
	}
}

func (o *OutputBasicListPaging) GetResult() map[string]interface{} {
	return map[string]interface{}{
		"code":        http.StatusOK,
		"result":      "OK",
		"list":        o.List,
		"count_total": o.CountTotal,
		"page":        o.Page,
		"limit":       o.Limit,
	}
}

func (o *OutputBasicListPaging) GetCode() int {
	return o.Code
}
