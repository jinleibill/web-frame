package web_frame

import (
	"fmt"
	"net/http"
	"strings"
)

type node struct {
	route string
	path  string

	children map[string]*node

	paramChild *node

	handler HandleFunc

	mdls []Middleware
}

func (n *node) childOrCreate(seg string) *node {
	if seg[0] == ':' {
		if n.paramChild == nil {
			n.paramChild = &node{
				path: seg,
			}
		} else {
			if n.paramChild.path != seg {
				panic("web 路由冲突，重复注册[:xxx]")
			}
		}

		return n.paramChild
	}

	if n.children == nil {
		n.children = map[string]*node{}
	}

	child, ok := n.children[seg]
	if !ok {
		child = &node{
			path: seg,
		}
		n.children[seg] = child
	}
	return child
}

func (n *node) childOf(seg string) (*node, bool, bool) {
	if n.children == nil {
		return n.paramChild, n.paramChild != nil, n.paramChild != nil
	}
	child, ok := n.children[seg]
	if !ok {
		return n.paramChild, n.paramChild != nil, n.paramChild != nil
	}
	return child, false, ok
}

type router struct {
	trees map[string]*node
}

func (r *router) add(method string, path string, handleFunc HandleFunc, mdls ...Middleware) {
	if path == "" {
		panic("web 路径不能为空")
	}

	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	if path[0] != '/' {
		panic("web 路径必须以 / 开头")
	}

	if path == "/" {
		if root.handler != nil {
			panic("web 路由冲突，重复注册[/]")
		}
		root.handler = handleFunc
		root.route = "/"
		return
	}

	if path[len(path)-1] == '/' {
		panic("web 路径不能以 / 结尾")
	}

	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if seg == "" {
			panic("web 路径不能有连续 /")
		}
		child := root.childOrCreate(seg)
		root = child
	}

	if root.handler != nil {
		panic(fmt.Sprintf("web 路由冲突，重复注册[%s]", path))
	}

	root.handler = handleFunc
	root.route = path
	root.mdls = mdls
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}

	segs := strings.Split(strings.Trim(path, "/"), "/")
	var pathParams map[string]string
	for _, seg := range segs {
		child, param, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		if param {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			pathParams[child.path[1:]] = seg
		}
		root = child
	}

	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

type routerGroup struct {
	*router

	name string
	mdls []Middleware
}

func (rg *routerGroup) Group(name string, mdls ...Middleware) *routerGroup {
	if len(name) > 0 {
		name = rg.name + "/" + name
	}

	return &routerGroup{
		router: rg.router,
		name:   name,
		mdls:   append(rg.mdls, mdls...),
	}
}

func (rg *routerGroup) addRoute(method string, path string, handleFunc HandleFunc) {
	rg.add(method, rg.name+path, handleFunc, rg.mdls...)
}

func (rg *routerGroup) Get(path string, handleFunc HandleFunc) {
	rg.addRoute(http.MethodGet, path, handleFunc)
}

func (rg *routerGroup) Post(path string, handleFunc HandleFunc) {
	rg.addRoute(http.MethodPost, path, handleFunc)
}

func (rg *routerGroup) Put(path string, handleFunc HandleFunc) {
	rg.addRoute(http.MethodPut, path, handleFunc)
}

func (rg *routerGroup) Delete(path string, handleFunc HandleFunc) {
	rg.addRoute(http.MethodDelete, path, handleFunc)
}

func (rg *routerGroup) Options(path string, handleFunc HandleFunc) {
	rg.addRoute(http.MethodOptions, path, handleFunc)
}

func newRouterGroup() *routerGroup {
	return &routerGroup{
		router: newRouter(),
		mdls:   make([]Middleware, 0),
	}
}
