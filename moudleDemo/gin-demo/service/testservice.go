package service

import (
	"worker/models"
)

type Test struct {
	Id      int    `json:"id"`
	Testcol string `json:"testcol"`
}

func (test *Test) Insert() (id int, err error) {
	var testModel models.Test
	testModel.Id = test.Id
	testModel.Testcol = test.Testcol
	id, err = testModel.Insert()
	return
}
