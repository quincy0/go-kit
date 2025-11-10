package httpc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quincy0/go-kit/rest/router"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	type Req struct {
		Name string `json:"name"`
		Age  uint64 `json:"age"`
	}
	type Reply struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Name string `json:"name"`
			Age  uint64 `json:"age"`
		}
	}

	r1 := Req{
		Name: "test",
		Age:  18,
	}

	reply1 := Reply{}

	rt := router.NewRouter()
	err := rt.Handle(http.MethodPost, "/test/post", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bBytes, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)

		r2 := Req{}
		json.Unmarshal(bBytes, &r2)
		assert.Equal(t, r1.Name, r2.Name)
		assert.Equal(t, r1.Age, r2.Age)

		reply2 := Reply{
			Code: 0,
			Msg:  "ok",
			Data: struct {
				Name string `json:"name"`
				Age  uint64 `json:"age"`
			}{
				Name: r1.Name,
				Age:  r1.Age,
			},
		}

		bytes, err := json.Marshal(reply2)
		assert.Nil(t, err)

		w.Write(bytes)
	}))

	assert.Nil(t, err)

	srv := httptest.NewServer(http.HandlerFunc(rt.ServeHTTP))
	defer srv.Close()

	ap, _ := NewApp("test-post", srv.URL)
	c1, err := ap.Client().WithData(r1).Post("test/post", &reply1)

	assert.Nil(t, err)
	assert.Equal(t, c1.HttpCode(), http.StatusOK)
	assert.Equal(t, 0, reply1.Code)
	assert.Equal(t, "ok", reply1.Msg)
	assert.Equal(t, r1.Name, reply1.Data.Name)
	assert.Equal(t, r1.Age, reply1.Data.Age)
}

func TestPostBodyIsNil(t *testing.T) {
	type Req struct {
		Name string `json:"name"`
		Age  uint64 `json:"age"`
	}

	r1 := Req{
		Name: "test",
		Age:  18,
	}

	rt := router.NewRouter()
	err := rt.Handle(http.MethodPost, "/test/post", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bBytes, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)

		r2 := Req{}
		json.Unmarshal(bBytes, &r2)
		assert.Equal(t, r1.Name, r2.Name)
		assert.Equal(t, r1.Age, r2.Age)
	}))

	assert.Nil(t, err)

	srv := httptest.NewServer(http.HandlerFunc(rt.ServeHTTP))
	defer srv.Close()

	ap, _ := NewApp("test-post", srv.URL)
	c1, err := ap.Client().WithData(r1).Post("test/post", nil)

	assert.Nil(t, err)
	assert.Equal(t, c1.HttpCode(), http.StatusOK)
}

func TestGet(t *testing.T) {
	q1 := map[string]string{
		"name": "test",
		"age":  "18",
	}

	type Reply struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Name string `json:"name"`
			Age  uint64 `json:"age"`
		}
	}

	reply1 := Reply{
		Code: 0,
		Msg:  "ok",
		Data: struct {
			Name string `json:"name"`
			Age  uint64 `json:"age"`
		}{
			Name: q1["name"],
			Age:  18,
		},
	}

	rt := router.NewRouter()
	err := rt.Handle(http.MethodGet, "/test/post", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		age := r.URL.Query().Get("age")

		assert.Equal(t, q1["name"], name)
		assert.Equal(t, q1["age"], age)

		retBytes, err := json.Marshal(reply1)
		assert.Nil(t, err)
		w.Write(retBytes)
	}))

	assert.Nil(t, err)

	srv := httptest.NewServer(http.HandlerFunc(rt.ServeHTTP))
	defer srv.Close()

	reply2 := Reply{}

	ap, _ := NewApp("test-post", srv.URL)
	c1, err := ap.Client().WithQuery(q1).Get("test/post", &reply2)

	assert.Nil(t, err)
	assert.Equal(t, c1.HttpCode(), http.StatusOK)
	assert.Equal(t, reply2.Msg, reply1.Msg)
	assert.Equal(t, reply2.Data.Name, reply1.Data.Name)
	assert.Equal(t, reply2.Data.Age, reply1.Data.Age)
	assert.Equal(t, reply2.Code, reply1.Code)
}

func TestWithHeader(t *testing.T) {

	var res interface{}

	ap, err := NewApp("stripo", "https://stripo.email/emailgeneration/v1")
	if err != nil {
		t.Fatal(err)
	}

	reply, err := ap.Client().WithHeader(map[string]string{
		"Stripo-Api-Auth": "eyJhbGciOiJIUzI1NiJ9.eyJzZWN1cml0eUNvbnRleHQiOiJ7XCJhcGlLZXlcIjpcImEwODA4MDg0LTdiYzMtNDYzNy05ZWVlLTJlNmI3MGYxZDYxNFwiLFwicHJvamVjdElkXCI6ODIwNzkyfSJ9.OXAgo8YGu-k9jiYoVFoebLTO3HXIuG18nQRYAM0UHhw",
	}).Get("emails", &res)
	t.Log(reply.HttpCode())

	ret, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(ret))
}
