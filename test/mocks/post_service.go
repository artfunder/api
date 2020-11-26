package mocks

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/artfunder/transport"
	"github.com/artfunder/transport/test"
	"github.com/artfunder/transport/test/samples"
)

type postsService struct {
	posts []test.Post
}

// PostsRouter is a transport router
//
// creates a new instance of the postsService on every call
// to refresh posts in memory
func PostsRouter(path string, method string) (transport.Endpoint, error) {
	s := newPostsService()

	if isRootPath(path) {
		return s.getRootEndpoint(method)
	}
	if isIDPath(path) {
		return s.getIDEndpoint(path, method)
	}
	if path == "/internal-error" {
		return s.internalErrorEndpoint(), nil
	}
	return nil, errors.New("Bad Request")
}

func newPostsService() *postsService {
	service := new(postsService)
	service.posts = make([]test.Post, len(samples.SamplePosts))
	copy(service.posts, samples.SamplePosts)
	return service
}

func isRootPath(path string) bool {
	matchesRoot, _ := regexp.MatchString("^/posts/?$", path)
	return matchesRoot
}

func isIDPath(path string) bool {
	matchesID, _ := regexp.MatchString(`^/posts/\d+/?$`, path)
	return matchesID
}

func (s postsService) getRootEndpoint(method string) (transport.Endpoint, error) {
	if method == http.MethodGet {
		return s.getAllPostsEndpoint(), nil
	}
	if method == http.MethodPost {
		return s.createPostEndpoint(), nil
	}
	return nil, transport.ErrorBadMethod
}

func (s postsService) getIDEndpoint(
	path string, method string,
) (transport.Endpoint, error) {

	id, err := getIDFromPath(path)
	if err != nil {
		return nil, err
	}

	return s.getIDEndpointFromMethod(id, method)
}

func (s postsService) getIDEndpointFromMethod(
	id int, method string,
) (transport.Endpoint, error) {

	switch method {
	case http.MethodGet:
		return s.getOnePostEndpoint(id), nil

	case http.MethodPatch:
		return s.updatePostEndpoint(id), nil

	case http.MethodDelete:
		return s.deletePostEndpoint(id), nil

	default:
		return nil, transport.ErrorBadMethod
	}
}

func getIDFromPath(path string) (int, error) {
	parts := strings.Split(path, "/")
	return strconv.Atoi(parts[len(parts)-1])
}

func (s postsService) getAllPostsEndpoint() transport.Endpoint {
	return transport.EndpointFunc(func(body []byte) ([]byte, error) {
		return json.Marshal(s.posts)
	})
}
func (s postsService) createPostEndpoint() transport.Endpoint {
	return transport.EndpointFunc(func(body []byte) ([]byte, error) {
		var post test.Post
		err := json.Unmarshal(body, &post)
		if err != nil {
			return nil, errors.New("Bad JSON")
		}
		post.ID = len(s.posts) + 1
		return json.Marshal(post)
	})
}
func (s postsService) getOnePostEndpoint(id int) transport.Endpoint {
	return transport.EndpointFunc(func(body []byte) ([]byte, error) {
		for _, post := range s.posts {
			if post.ID == id {
				return json.Marshal(post)
			}
		}
		return nil, transport.ErrorNotFound
	})
}
func (s postsService) updatePostEndpoint(id int) transport.Endpoint {
	return transport.EndpointFunc(func(body []byte) ([]byte, error) {
		for i, post := range s.posts {
			if post.ID == id {
				err := json.Unmarshal(body, &s.posts[i])
				if err != nil {
					return nil, errors.New("Bad JSON")
				}
				s.posts[i].ID = post.ID
				return json.Marshal(s.posts[i])
			}
		}
		return nil, transport.ErrorNotFound
	})
}
func (s postsService) deletePostEndpoint(id int) transport.Endpoint {
	return transport.EndpointFunc(func(body []byte) ([]byte, error) {
		for i, post := range s.posts {
			if post.ID == id {
				s.posts = append(s.posts[:i], s.posts[i+1:]...)
				return json.Marshal(post)
			}
		}
		return nil, transport.ErrorNotFound
	})
}

func (s postsService) internalErrorEndpoint() transport.Endpoint {
	return transport.EndpointFunc(func(body []byte) ([]byte, error) {
		return []byte("eggplant"), nil
	})
}
