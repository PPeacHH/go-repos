package main

import "testing"

func TestPlateHandler_GetAccessToken(t *testing.T) {
	handler := PlateHandler{}
	a := "----"
	b := "----"
	if _, err :=handler.GetAccessToken(a, b); err != nil {
		t.Errorf("Handler_GetAccessToken has happened some error:%v", err)
	}
}
