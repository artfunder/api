package transport_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/artfunder/structs"
	"github.com/artfunder/transport"
	"github.com/artfunder/transport/test"
	"github.com/artfunder/transport/test/mocks"
	"github.com/artfunder/transport/test/samples"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
)

type FastHTTPTest struct {
	name         string
	router       transport.Router
	action       func() (*http.Response, error)
	decodeResp   func([]byte) (interface{}, error)
	expectedResp interface{}
}

var fastHTTPTests = []FastHTTPTest{
	{
		name:       "Create Service",
		router:     mocks.PostsRouter,
		decodeResp: decodeToPost,
		action: func() (*http.Response, error) {
			body, _ := json.Marshal(map[string]interface{}{
				"title":   "foo",
				"content": "bar",
				"likes":   2,
			})

			return http.Post(
				"http://localhost:10000/posts",
				"application/json",
				bytes.NewBuffer(body),
			)
		},
		expectedResp: test.Post{
			ID:      4,
			Title:   "foo",
			Content: "bar",
			Likes:   2,
		},
	},
	{
		name:   "Find Service",
		router: mocks.PostsRouter,
		action: func() (*http.Response, error) {
			return http.Get("http://localhost:10000/posts")
		},
		decodeResp: func(b []byte) (interface{}, error) {
			var posts []test.Post
			err := json.Unmarshal(b, &posts)
			return posts, err
		},
		expectedResp: samples.SamplePosts,
	},
	{
		name:   "Get-One Service",
		router: mocks.PostsRouter,
		action: func() (*http.Response, error) {
			return http.Get("http://localhost:10000/posts/1")
		},
		decodeResp: decodeToPost,
		expectedResp: test.Post{
			ID:      1,
			Title:   "Post 1",
			Content: "Content of Post 1",
			Likes:   3,
		},
	},
	{
		name:   "Update Service",
		router: mocks.PostsRouter,
		action: func() (*http.Response, error) {
			body, _ := json.Marshal(map[string]interface{}{
				"title":   "foo",
				"content": "bar",
			})
			req, _ := http.NewRequest(
				"PATCH",
				"http://localhost:10000/posts/2",
				bytes.NewBuffer(body),
			)
			return http.DefaultClient.Do(req)
		},
		decodeResp: decodeToPost,
		expectedResp: test.Post{
			ID:      2,
			Title:   "foo",
			Content: "bar",
			Likes:   2,
		},
	},
	{
		name:   "Delete Service",
		router: mocks.PostsRouter,
		action: func() (*http.Response, error) {
			req, _ := http.NewRequest("DELETE", "http://localhost:10000/posts/2", nil)
			return http.DefaultClient.Do(req)
		},
		decodeResp: decodeToPost,
		expectedResp: test.Post{
			ID:      2,
			Title:   "Post 2",
			Content: "Content of Post 2",
			Likes:   2,
		},
	},
	{
		name:   "Bad URL",
		router: mocks.PostsRouter,
		action: func() (*http.Response, error) {
			return http.Get("http://localhost:10000/badpath")
		},
		decodeResp: decodeToError,
		expectedResp: structs.Error{
			Message: "Bad Request",
		},
	},
	{
		name:   "Bad JSON",
		router: mocks.PostsRouter,
		action: func() (*http.Response, error) {
			return http.Post(
				"http://localhost:10000/posts",
				"application/json",
				bytes.NewBuffer([]byte("eggplant")),
			)
		},
		decodeResp: decodeToError,
		expectedResp: structs.Error{
			Message: "Bad JSON",
		},
	},
	{
		name:   "Out of Bounds",
		router: mocks.PostsRouter,
		action: func() (*http.Response, error) {
			return http.Get("http://localhost:10000/posts/5")
		},
		decodeResp: decodeToError,
		expectedResp: structs.Error{
			Message: "Object Not Found",
		},
	},
	{
		name:   "Bad Method",
		router: mocks.PostsRouter,
		action: func() (*http.Response, error) {
			req, _ := http.NewRequest(http.MethodPut, "http://localhost:10000/posts", nil)
			return http.DefaultClient.Do(req)
		},
		decodeResp: decodeToError,
		expectedResp: structs.Error{
			Message: "Method Not Allowed",
		},
	},
	{
		name:   "Internal Error",
		router: mocks.PostsRouter,
		action: func() (*http.Response, error) {
			return http.Get("http://localhost:10000/internal-error")
		},
		decodeResp: decodeToError,
		expectedResp: structs.Error{
			Message: "Internal Error",
		},
	},
}

func TestFastHTTP(t *testing.T) {
	fhTransport, ln := startServer(t)
	defer ln.Close()

	for _, test := range fastHTTPTests {
		t.Run(test.name, func(t *testing.T) {
			fhTransport.Route(test.router)

			resp, err := test.action()
			if err != nil {
				t.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			actual, err := test.decodeResp(body)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.expectedResp, actual); diff != "" {
				t.Fatal(diff, "Actual string:", string(body))
			}
		})
	}
}

func startServer(t *testing.T) (*transport.FastHTTPTransport, io.Closer) {
	fhTransport := new(transport.FastHTTPTransport)
	ln, err := net.Listen("tcp", "localhost:10000")
	if err != nil {
		t.Fatal("Couldn't start server")
	}
	go fasthttp.Serve(ln, fhTransport.HandleFastHTTP)
	return fhTransport, ln
}

func decodeToPost(b []byte) (interface{}, error) {
	var p test.Post
	err := json.Unmarshal(b, &p)
	return p, err
}
func decodeToError(b []byte) (interface{}, error) {
	var e structs.Error
	err := json.Unmarshal(b, &e)
	return e, err
}
