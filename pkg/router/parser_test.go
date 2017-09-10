package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type registered struct {
	path string
	h    handle
}

type expected struct {
	request    string
	data       string
	paramCount int
	params     []param
}

var setOfRegistered = []registered{
	{
		"/hello/John",
		func(c Control) {
			c.Write("Hello from static path")
		},
	},
	{
		"/hello/:name",
		func(c Control) {
			c.Write("Hello " + c.Query(":name"))
		},
	},
	{
		"/:h/:n",
		func(c Control) {
			c.Write(c.Query(":n") + " from " + c.Query(":h"))
		},
	},
	{
		"/products/book/orders/:id",
		func(c Control) {
			c.Write("Product: book order# " + c.Query(":id"))
		},
	},
	{
		"/products/:name/orders/:id",
		func(c Control) {
			c.Write("Product: " + c.Query(":name") + " order# " + c.Query(":id"))
		},
	},
	{
		"/products/:name/:order/:id",
		func(c Control) {
			c.Write("Product: " + c.Query(":name") + " # " + c.Query(":id"))
		},
	},
	{
		"/:product/:name/:order/:id",
		func(c Control) {
			c.Write(c.Query(":product") + " " + c.Query(":name") + " " + c.Query(":order") + " # " + c.Query(":id"))
		},
	},
	{
		"/static/*",
		func(c Control) {
			c.Write("Hello from star static path")
		},
	},
	{
		"/files/:dir/*",
		func(c Control) {
			c.Write(c.Query(":dir"))
		},
	},
}

var setOfExpected = []expected{
	{
		"/hello/John",
		"Hello from static path",
		0,
		[]param{},
	},
	{
		"/hello/Jane",
		"Hello Jane",
		1,
		[]param{
			{":name", "Jane"},
		},
	},
	{
		"/hell/jack",
		"jack from hell",
		2,
		[]param{
			{":h", "hell"},
			{":n", "jack"},
		},
	},
	{
		"/products/book/orders/12",
		"Product: book order# 12",
		1,
		[]param{
			{":id", "12"},
		},
	},
	{
		"/products/table/orders/23",
		"Product: table order# 23",
		2,
		[]param{
			{":name", "table"},
			{":id", "23"},
		},
	},
	{
		"/products/pen/orders/11",
		"Product: pen order# 11",
		2,
		[]param{
			{":name", "pen"},
			{":id", "11"},
		},
	},
	{
		"/products/pen/order/10",
		"Product: pen # 10",
		3,
		[]param{
			{":name", "pen"},
			{":order", "order"},
			{":id", "10"},
		},
	},
	{
		"/product/pen/order/10",
		"product pen order # 10",
		4,
		[]param{
			{":product", "product"},
			{":name", "pen"},
			{":order", "order"},
			{":id", "10"},
		},
	},
	{
		"/static/greetings/something",
		"Hello from star static path",
		0,
		[]param{},
	},
	{
		"/files/css/style.css",
		"css",
		1,
		[]param{
			{":dir", "css"},
		},
	},
	{
		"/files/js/app.js",
		"js",
		1,
		[]param{
			{":dir", "js"},
		},
	},
}

func TestParserRegisterGet(t *testing.T) {
	p := newParser()
	for _, request := range setOfRegistered {
		p.register(request.path, request.h)
	}
	for _, exp := range setOfExpected {
		h, params, ok := p.get(exp.request)
		if !ok {
			t.Error("Error: get data for path", exp.request)
		}
		if len(params) != exp.paramCount {
			t.Error("Expected length of param", exp.paramCount, "got", len(params))
		}
		trw := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Error("Error creating new request")
		}
		c := newControl(trw, req)
		for _, item := range params {
			c.Param(item.key, item.value)
		}
		h(c)
		if trw.Body.String() != exp.data {
			t.Error("Expected", exp.data, "got", trw.Body.String())
		}
	}
}

func TestParserSplit(t *testing.T) {
	path := []string{
		"/api/v1/module",
		"/api//v1/module/",
		"/module///name//",
		"module//:name",
		"/:param1/:param2/",
		strings.Repeat("/A", 300),
	}
	expected := [][]string{
		{"api", "v1", "module"},
		{"api", "v1", "module"},
		{"module", "name"},
		{"module", ":name"},
		{":param1", ":param2"},
	}

	if part, ok := split("   "); ok {
		if len(part) != 0 {
			t.Error("Error: split data for path '/'", part)
		}
	} else {
		t.Error("Error: split data for path '/'")
	}

	if part, ok := split("///"); ok {
		if len(part) != 0 {
			t.Error("Error: split data for path '/'", part)
		}
	} else {
		t.Error("Error: split data for path '/'")
	}

	if part, ok := split("  /  //  "); ok {
		if len(part) != 0 {
			t.Error("Error: split data for path '/'", part)
		}
	} else {
		t.Error("Error: split data for path '/'")
	}

	for idx, p := range path {
		parts, ok := split(p)
		if !ok {
			if strings.HasPrefix(p, "/A/A/A") {
				parser := newParser()
				result := parser.register(p, func(Control) {})
				if result {
					t.Error("Expected false result, got", result)
				}
				continue
			}
			t.Error("Error: split data for path", p)
		}
		for i, part := range parts {
			if expected[idx][i] != part {
				t.Error("Expected", expected[idx][i], "got", part)
			}
		}
	}
}

func TestGetRoutes(t *testing.T) {
	for _, request := range setOfRegistered {
		p := newParser()
		p.register(request.path, request.h)
		routes := p.routes()
		if len(routes) != 1 {
			t.Error("Expected 1 route, got", len(routes))
		}
		if request.path != routes[0] {
			t.Error("Expected", request.path, "got", routes[0])
		}
	}
}

func TestRegisterAsterisk(t *testing.T) {
	data := "Any path is ok"
	p := newParser()
	p.register("*", func(c Control) {
		c.Write(data)
	})
	path := "/any/path/is/ok"
	h, params, ok := p.get(path)
	if !ok {
		t.Error("Error: get data for path", path)
	}
	trw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Error("Error creating new request")
	}
	c := newControl(trw, req)
	for _, item := range params {
		c.Param(item.key, item.value)
	}
	h(c)
	if trw.Body.String() != data {
		t.Error("Expected", data, "got", trw.Body.String())
	}
}
