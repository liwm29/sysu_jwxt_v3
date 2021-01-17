package jwxtClient

var (
	_campus_id = map[string]string{
		"东校园":  "5063559",
		"北校园":  "5062202",
		"南校园":  "5062201",
		"深圳校区": "333291143",
		"珠海校区": "5062203",
	}
	_selectedType = map[string]string{
		"本专业":  "1",
		"校级公选": "4",
		"跨专业":  "2",
	}
	_selectedCate = map[string]string{
		"专必":     "11",
		"专选":     "21",
		"院内公选":   "30",
		"公必(体育)": "10",
		"公必(大英)": "10",
		"公必(其他)": "10",
	}
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
