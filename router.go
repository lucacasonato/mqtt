package mqtt

import (
	"strings"
	"sync"
)

type router struct {
	routes []route
	lock   sync.RWMutex
}

func newRouter() *router {
	return &router{routes: []route{}, lock: sync.RWMutex{}}
}

type route struct {
	topic   string
	handler MessageHandler
}

func match(route []string, topic []string) bool {
	if len(route) == 0 {
		return len(topic) == 0
	}

	if len(topic) == 0 {
		return route[0] == "#"
	}

	if route[0] == "#" {
		return true
	}

	if (route[0] == "+") || (route[0] == topic[0]) {
		return match(route[1:], topic[1:])
	}
	return false
}

func routeIncludesTopic(route, topic string) bool {
	return match(routeSplit(route), strings.Split(topic, "/"))
}

func routeSplit(route string) []string {
	var result []string
	if strings.HasPrefix(route, "$share") {
		result = strings.Split(route, "/")[2:]
	} else {
		result = strings.Split(route, "/")
	}
	return result
}

func (r *route) match(message *Message) bool {
	return r.topic == message.Topic() || routeIncludesTopic(r.topic, message.Topic())
}

func (r *router) addRoute(topic string, handler MessageHandler) {
	if handler != nil {
		r.lock.Lock()
		r.routes = append(r.routes, route{topic: topic, handler: handler})
		r.lock.Unlock()
	}
}

func (r *router) match(message *Message) []route {
	routes := []route{}
	r.lock.RLock()
	for _, route := range r.routes {
		if route.match(message) {
			routes = append(routes, route)
		}
	}
	r.lock.RUnlock()
	return routes
}
