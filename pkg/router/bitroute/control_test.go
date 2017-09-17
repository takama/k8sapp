package bitroute

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type prm struct {
	Key, Value string
}

var params = []prm{
	{"name", "John"},
	{"age", "32"},
	{"gender", "M"},
}

var testParamsData = `[{"Key":"name","Value":"John"},{"Key":"age","Value":"32"},{"Key":"gender","Value":"M"}]`
var testParamGzipData = []byte{
	31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 138, 174, 86, 242, 78, 173, 84, 178, 82, 202,
	75, 204, 77, 85, 210, 81, 10, 75, 204, 41, 77, 85, 178, 82, 242, 202, 207, 200,
	83, 170, 213, 129, 201, 38, 166, 35, 75, 26, 27, 33, 73, 165, 167, 230, 165, 164,
	22, 33, 201, 250, 42, 213, 198, 2, 2, 0, 0, 255, 255, 196, 73, 247, 37, 87, 0, 0, 0,
}

func TestParamsQueryGet(t *testing.T) {

	c := new(control)
	for _, param := range params {
		c.Param(param.Key, param.Value)
	}
	for _, param := range params {
		value := c.Query(param.Key)
		if value != param.Value {
			t.Error("Expected for", param.Key, ":", param.Value, ", got", value)
		}
	}
}

func TestWriterHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "hello/:name", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	c := NewControl(trw, req)
	request := c.Request()
	if request != req {
		t.Error("Expected", req.URL.String(), "got", request.URL.String())
	}
	trw.Header().Add("Test", "TestValue")
	c = NewControl(trw, req)
	expected := trw.Header().Get("Test")
	value := c.Header().Get("Test")
	if value != expected {
		t.Error("Expected", expected, "got", value)
	}
}

func TestWriterCode(t *testing.T) {
	c := new(control)
	// code transcends, must be less than 600
	c.Code(777)
	if c.code != 0 {
		t.Error("Expected code", "0", "got", c.code)
	}
	c.Code(404)
	if c.code != 404 {
		t.Error("Expected code", "404", "got", c.code)
	}
}

func TestGetCode(t *testing.T) {
	c := new(control)
	c.Code(http.StatusOK)
	code := c.GetCode()
	if code != http.StatusOK {
		t.Error("Expected code", http.StatusText(http.StatusOK), "got", code)
	}
}

func TestWrite(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	c := NewControl(trw, req)
	c.Write("Hello")
	if trw.Body.String() != "Hello" {
		t.Error("Expected", "Hello", "got", trw.Body.String())
	}
	contentType := trw.Header().Get("Content-type")
	expected := "text/plain; charset=utf-8"
	if contentType != expected {
		t.Error("Expected", expected, "got", contentType)
	}
	trw = httptest.NewRecorder()
	c = NewControl(trw, req)
	c.Code(http.StatusOK)
	c.Write(params)
	if trw.Body.String() != testParamsData {
		t.Error("Expected", testParamsData, "got", trw.Body.String())
	}
	contentType = trw.Header().Get("Content-type")
	expected = "application/json"
	if contentType != expected {
		t.Error("Expected", expected, "got", contentType)
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	trw = httptest.NewRecorder()
	c = NewControl(trw, req)
	c.Code(http.StatusAccepted)
	c.Write(params)
	if trw.Body.String() != string(testParamGzipData) {
		t.Error("Expected", testParamGzipData, "got", trw.Body.String())
	}
	contentEncoding := trw.Header().Get("Content-Encoding")
	expected = "gzip"
	if contentEncoding != expected {
		t.Error("Expected", expected, "got", contentEncoding)
	}
	trw = httptest.NewRecorder()
	c = NewControl(trw, req)
	c.Write(func() {})
	if trw.Code != http.StatusInternalServerError {
		t.Error("Expected", http.StatusInternalServerError, "got", trw.Code)
	}
	expected = "application/json"
	if contentType != expected {
		t.Error("Expected", expected, "got", contentType)
	}
}
