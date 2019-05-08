package models

const (
	RECODE_OK        = "0"
	RECODE_NODATA    = "4002"
	RECODE_DATAERR   = "4004"
	RECODE_PARAMERR  = "4103"
	RECODE_UNKNOWERR = "4501"
)

var recodeText = map[string]string{
	RECODE_OK:        "成功",
	RECODE_NODATA:    "无数据",
	RECODE_DATAERR:   "数据错误",
	RECODE_PARAMERR:  "参数错误",
	RECODE_UNKNOWERR: "未知错误",
}

func RecodeText(code string) string {
	str, ok := recodeText[code]
	if ok {
		return str
	}
	return recodeText[RECODE_UNKNOWERR]
}
