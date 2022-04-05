package database

// Variables 请求响应变量表结构体
type Variables struct {
	Id                string `db:"id"`
	TemplateName      string `db:"template_name"`
	TemplateVersion   int64  `db:"template_version"`
	RequestHeader     string `db:"request_header"`
	RequestPathParam  string `db:"request_path_params"`
	RequestQueryParam string `db:"request_query_params"`
	RequestBodyVar    string `db:"request_body_var"`
	ResponseHeader    string `db:"response_header"`
	ResponseBodyVar   string `db:"response_body_var"`
}

func (vars *Variables) Insert() (int64, int64, error) {
	return 0, 0, nil
}
