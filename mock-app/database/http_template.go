package database

//HttpTemplate http请求模板表
type HttpTemplate struct {
	Id                   string `db:"id"`
	Name                 string `db:"name"`
	DisplayName          string `db:"display_name"`
	Version              int64  `db:"version"`
	RequestBodyTemplate  string `db:"request_body_template"`
	ResponseBodyTemplate string `db:"response_body_template"`
}
