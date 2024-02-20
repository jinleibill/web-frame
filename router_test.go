package web_frame

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_addRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail/:id",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodGet,
			path:   "/items",
		},
		{
			method: http.MethodGet,
			path:   "/items/:item/comments/:comment",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.add(route.method, route.path, mockHandler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": {
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
								paramChild: &node{
									path:    ":id",
									handler: mockHandler,
								},
							},
						},
					},
					"items": {
						path:    "items",
						handler: mockHandler,
						paramChild: &node{
							path: ":item",
							children: map[string]*node{
								"comments": {
									path: "comments",
									paramChild: &node{
										path:    ":comment",
										handler: mockHandler,
									},
								},
							},
						},
					},
				},
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"create": {
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login": {
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}

	msg, ok := r.equal(wantRouter)
	assert.True(t, ok, msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.add(http.MethodGet, "", mockHandler)
	}, "web 路径不能为空")
	assert.Panicsf(t, func() {
		r.add(http.MethodGet, "login", mockHandler)
	}, "web 路径必须以 / 开头")
	assert.Panicsf(t, func() {
		r.add(http.MethodGet, "/login/", mockHandler)
	}, "web 路径不能以 / 结尾")
	assert.Panicsf(t, func() {
		r.add(http.MethodGet, "//login", mockHandler)
	}, "web 路径不能有连续 /")

	r = newRouter()
	assert.Panicsf(t, func() {
		r.add(http.MethodGet, "/", mockHandler)
		r.add(http.MethodGet, "/", mockHandler)
	}, "web 路由冲突，重复注册[/]")

	r = newRouter()
	assert.Panicsf(t, func() {
		r.add(http.MethodGet, "/login", mockHandler)
		r.add(http.MethodGet, "/login", mockHandler)
	}, "web 路由冲突，重复注册[/login]")

	r = newRouter()
	assert.Panicsf(t, func() {
		r.add(http.MethodGet, "/banner/:id", mockHandler)
		r.add(http.MethodGet, "/banner/:name", mockHandler)
	}, "web 路由冲突，重复注册[/banner/:xxx]")
}

func (r *router) equal(tr *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := tr.trees[k]
		if !ok {
			return fmt.Sprintf("http method 方法不匹配"), false
		}
		msg, ok := v.equal(dst)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

func (n *node) equal(tn *node) (string, bool) {
	if n.path != tn.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.children) != len(tn.children) {
		return fmt.Sprintf("子节点数量不相等"), false
	}
	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(tn.paramChild)
		if !ok {
			return msg, false
		}
	}
	nHandler := reflect.ValueOf(n.handler)
	tnHandler := reflect.ValueOf(tn.handler)
	if nHandler != tnHandler {
		return fmt.Sprintf("handler 不相等"), false
	}

	for path, c := range n.children {
		dst, ok := tn.children[path]
		if !ok {
			return fmt.Sprintf("子节点不存在: %s", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodDelete,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/banner",
		},
		{
			method: http.MethodGet,
			path:   "/banner/:id",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	for _, router := range testRoutes {
		r.add(router.method, router.path, mockHandler)
	}

	testCases := []struct {
		name string

		method string
		path   string

		wantFound bool
		info      *matchInfo
	}{
		{
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "/",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "method not found",
			method:    http.MethodPut,
			path:      "/order/create",
			wantFound: false,
		},
		{
			name:      "path not found",
			method:    http.MethodGet,
			path:      "/user",
			wantFound: false,
		},
		{
			name:      "banner",
			method:    http.MethodGet,
			path:      "/banner",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "banner",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "banner id",
			method:    http.MethodGet,
			path:      "/banner/:id",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    ":id",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"id": ":id",
				},
			},
		},
		{
			name:      "order create",
			method:    http.MethodPost,
			path:      "/order/create",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "create",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "order not handler",
			method:    http.MethodPost,
			path:      "/order",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path: "order",
					children: map[string]*node{
						"create": {
							path:    "create",
							handler: mockHandler,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.info.pathParams, info.pathParams)
			meg, ok := tc.info.n.equal(info.n)
			assert.True(t, ok, meg)
		})
	}

}

func TestRouterGroup_Group(t *testing.T) {
	type childGroup struct {
		name string

		wantName string
	}

	testRouteGroups := []struct {
		name             string
		childRouteGroups []childGroup

		wantName string
	}{
		{
			name: "v1",
			childRouteGroups: []childGroup{
				{name: "admins", wantName: "/v1/admins"},
				{name: "users", wantName: "/v1/users"},
			},
			wantName: "/v1",
		},
	}

	rg := newRouterGroup()
	for _, group := range testRouteGroups {
		fGroup := rg.Group(group.name)
		assert.Equal(t, fGroup.name, group.wantName, group.name)
		for _, childGroup := range group.childRouteGroups {
			cGroup := fGroup.Group(childGroup.name)
			assert.Equal(t, cGroup.name, childGroup.wantName, childGroup.name)
		}
	}
}
