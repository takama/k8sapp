package bitroute

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/takama/k8sapp/pkg/router"
)

func getRouterForTesting() *bitroute {
	return &bitroute{
		handlers: make(map[string]*parser),
	}
}

func TestNewRouter(t *testing.T) {
	r := New()
	if r == nil {
		t.Error("Expected new router, got nil")
	}
	err := r.Listen("$")
	if err == nil {
		t.Error("Expected error if used incorrect host and port")
	}
}

func TestRouterGetRootStatic(t *testing.T) {
	r := getRouterForTesting()
	// Registers GET handler for root static path
	r.GET("/", func(c router.Control) {
		c.Write("Root")
	})
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	if trw.Body.String() != "Root" {
		t.Error("Expected", "Root", "got", trw.Body.String())
	}
}

func TestRouterGetStatic(t *testing.T) {
	r := getRouterForTesting()
	// Registers GET handler for static path
	r.GET("/hello", func(c router.Control) {
		c.Write("Hello")
	})
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	if trw.Body.String() != "Hello" {
		t.Error("Expected", "Hello", "got", trw.Body.String())
	}
}

func TestRouterGetParameter(t *testing.T) {
	r := getRouterForTesting()
	// Registers GET handler with parameter
	r.GET("/hello/:name", func(c router.Control) {
		c.Write("Hello " + c.Query(":name"))
	})
	req, err := http.NewRequest("GET", "/hello/John", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	if trw.Body.String() != "Hello John" {
		t.Error("Expected", "Hello John", "got", trw.Body.String())
	}
}

func TestRouterGetParameterFromClassicUrl(t *testing.T) {
	r := getRouterForTesting()
	// Registers GET handler with two parameters
	r.GET("/users/:name", func(c router.Control) {
		c.Write("Users: " + c.Query(":name") + " " + c.Query("name"))
	})
	req, err := http.NewRequest("GET", "/users/Jane/?name=Joe", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	if trw.Body.String() != "Users: Jane Joe" {
		t.Error("Expected", "Users: Jane Joe", "got", trw.Body.String())
	}
}

func TestRouterPostJSONData(t *testing.T) {
	r := getRouterForTesting()
	// Registers POST handler
	r.POST("/users", func(c router.Control) {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			t.Error(err)
		}
		var values map[string]string
		if err := json.Unmarshal(body, &values); err != nil {
			t.Error(err)
		}
		c.Write("User: " + values["name"])
	})
	req, err := http.NewRequest("POST", "/users/", strings.NewReader(`{"name": "Tom"}`))
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	if trw.Body.String() != "User: Tom" {
		t.Error("Expected", "User: Tom", "got", trw.Body.String())
	}
}

func TestRouterPutJSONData(t *testing.T) {
	r := getRouterForTesting()
	// Registers PUT handler
	r.PUT("/users", func(c router.Control) {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			t.Error(err)
		}
		var values map[string]string
		if err := json.Unmarshal(body, &values); err != nil {
			t.Error(err)
		}
		c.Write("Users: " + values["name1"] + " " + values["name2"])
	})
	req, err := http.NewRequest("PUT", "/users/", strings.NewReader(`{"name1": "user1", "name2": "user2"}`))
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	if trw.Body.String() != "Users: user1 user2" {
		t.Error("Expected", "Users: user1 user2", "got", trw.Body.String())
	}
}

func TestRouterDelete(t *testing.T) {
	r := getRouterForTesting()
	// Registers DELETE handler
	r.DELETE("/users", func(c router.Control) {
		c.Write("Users deleted")
	})
	req, err := http.NewRequest("DELETE", "/users/", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	if trw.Body.String() != "Users deleted" {
		t.Error("Expected", "Users deleted", "got", trw.Body.String())
	}
}

func TestRouterHead(t *testing.T) {
	r := getRouterForTesting()
	// Registers HEAD handler
	r.HEAD("/command", func(c router.Control) {
		c.Header().Add("test", "value")
	})
	req, err := http.NewRequest("HEAD", "/command/", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	result := trw.Header().Get("test")
	if result != "value" {
		t.Error("Expected value", "got", result)
	}
}

func TestRouterOptions(t *testing.T) {
	r := getRouterForTesting()
	// Registers OPTIONS handler
	r.OPTIONS("/option", func(c router.Control) {
		c.Code(http.StatusOK)
	})
	req, err := http.NewRequest("OPTIONS", "/option/", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	result := trw.Code
	if result != http.StatusOK {
		t.Error("Expected", http.StatusOK, "got", result)
	}
}

func TestRouterPatch(t *testing.T) {
	r := getRouterForTesting()
	// Registers PATCH handler
	r.PATCH("/patch", func(c router.Control) {
		c.Code(http.StatusOK)
	})
	req, err := http.NewRequest("PATCH", "/patch/", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	result := trw.Code
	if result != http.StatusOK {
		t.Error("Expected", http.StatusOK, "got", result)
	}
}

func TestRouterUseOptionsReplies(t *testing.T) {
	r := getRouterForTesting()
	path := "/options"
	r.GET(path, func(c router.Control) {
		c.Code(http.StatusOK)
	})
	r.UseOptionsReplies(true)
	req, err := http.NewRequest("OPTIONS", path, nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	code := trw.Code
	if code != http.StatusOK {
		t.Error("Expected", http.StatusOK, "got", code)
	}
	header := trw.Header().Get("Allow")
	expected := "GET"
	if header != expected {
		t.Error("Expected", expected, "got", header)
	}
}

func TestRouterNotFound(t *testing.T) {
	r := getRouterForTesting()
	// Registers GET handler
	r.GET("/found", func(c router.Control) {
		c.Code(http.StatusOK)
	})
	req, err := http.NewRequest("GET", "/not-found/", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	result := trw.Code
	if result != http.StatusNotFound {
		t.Error("Expected", http.StatusNotFound, "got", result)
	}
}

func TestRouterAllowedMethods(t *testing.T) {
	r := getRouterForTesting()
	// Registers GET handler
	path := "/allowed"
	r.GET(path, func(c router.Control) {
		c.Code(http.StatusOK)
	})
	// Registers PUT handler
	r.PUT(path, func(c router.Control) {
		c.Code(http.StatusAccepted)
	})
	result := r.allowedMethods(path)
	for _, method := range []string{"GET", "PUT"} {
		var exists bool
		for _, allowed := range result {
			if method == allowed {
				exists = true
			}
		}
		if !exists {
			t.Error("Allowed method(s) not found in", result)
		}
	}
	for _, method := range []string{"POST", "DELETE", "HEAD", "OPTIONS", "PATCH"} {
		var exists bool
		for _, allowed := range result {
			if method == allowed {
				exists = true
			}
		}
		if exists {
			t.Error("Not allowed method(s) found in", result)
		}
	}
}

func TestRouterNotAllowed(t *testing.T) {
	r := getRouterForTesting()
	// Registers GET handler
	path := "/allowed"
	message := http.StatusText(http.StatusMethodNotAllowed) + "\n"
	r.GET(path, func(c router.Control) {
		c.Code(http.StatusOK)
	})
	// Registers PUT handler
	r.PUT(path, func(c router.Control) {
		c.Code(http.StatusAccepted)
	})
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	result := trw.Code
	if result != http.StatusMethodNotAllowed {
		t.Error("Expected", http.StatusMethodNotAllowed, "got", result)
	}
	header := trw.Header().Get("Allow")
	expected1 := "GET, PUT"
	expected2 := "PUT, GET"
	if header != expected1 && header != expected2 {
		t.Error("Expected", expected1, "or", expected2, "got", header)
	}
	if trw.Body.String() != message {
		t.Error("Expected", message, "got", trw.Body.String())
	}
}

func TestRouterSetupNotAllowedHandler(t *testing.T) {
	r := getRouterForTesting()
	message := http.StatusText(http.StatusForbidden)
	path := "/not/allowed"
	r.GET(path, func(c router.Control) {
		c.Code(http.StatusOK)
	})
	r.SetupNotAllowedHandler(func(c router.Control) {
		c.Code(http.StatusForbidden)
		c.Write(message)
	})
	req, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	code := trw.Code
	if code != http.StatusForbidden {
		t.Error("Expected", http.StatusForbidden, "got", code)
	}
	header := trw.Header().Get("Allow")
	expected := "GET"
	if header != expected {
		t.Error("Expected", expected, "got", header)
	}
	if trw.Body.String() != message {
		t.Error("Expected", message, "got", trw.Body.String())
	}
}

func TestRouterSetupNotFound(t *testing.T) {
	r := getRouterForTesting()
	message := http.StatusText(http.StatusForbidden)
	r.SetupNotFoundHandler(func(c router.Control) {
		c.Code(http.StatusForbidden)
		c.Write(message)
	})
	req, err := http.NewRequest("GET", "/not/found", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	result := trw.Code
	if result != http.StatusForbidden {
		t.Error("Expected", http.StatusForbidden, "got", result)
	}
	if trw.Body.String() != message {
		t.Error("Expected", message, "got", trw.Body.String())
	}
}

func TestRouterRecoveryHandler(t *testing.T) {
	r := getRouterForTesting()
	message := http.StatusText(http.StatusServiceUnavailable)
	path := "/recovery"
	r.GET(path, func(c router.Control) {
		panic("test")
	})
	r.SetupRecoveryHandler(func(c router.Control) {
		c.Code(http.StatusServiceUnavailable)
		c.Write(message)
	})
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	result := trw.Code
	if result != http.StatusServiceUnavailable {
		t.Error("Expected", http.StatusForbidden, "got", result)
	}
	if trw.Body.String() != message {
		t.Error("Expected", message, "got", trw.Body.String())
	}
}

func TestRouterMiddleware(t *testing.T) {
	r := getRouterForTesting()
	message := http.StatusText(http.StatusOK)
	path := "/middleware"
	r.GET(path, func(c router.Control) {
		c.Code(http.StatusOK)
		c.Write(message)
	})
	r.SetupMiddleware(func(f func(router.Control)) func(router.Control) {
		return func(c router.Control) {
			headers := c.Request().Header.Get("Access-Control-Request-Headers")
			if headers != "" {
				c.Header().Set("Access-Control-Allow-Headers", "content-type")
			}
			f(c)
		}
	})
	req, err := http.NewRequest("GET", path, nil)
	req.Header.Set("Access-Control-Request-Headers", "All")
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	result := trw.Code
	if result != http.StatusOK {
		t.Error("Expected", http.StatusOK, "got", result)
	}
	header := trw.Header().Get("Access-Control-Allow-Headers")
	expected := "content-type"
	if header != expected {
		t.Error("Expected", expected, "got", header)
	}
	if trw.Body.String() != message {
		t.Error("Expected", message, "got", trw.Body.String())
	}
}
