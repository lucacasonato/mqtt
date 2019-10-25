package mqtt

import (
	"strings"
	"sync"

	"github.com/google/uuid"
)

type router struct {
	routes []Route
	lock   sync.RWMutex
}

func newRouter() *router {
	return &router{routes: []Route{}, lock: sync.RWMutex{}}
}

type Route struct {
	router  *router
	id      string
	topic   string
	handler MessageHandler
}

func newRoute(router *router, topic string, handler MessageHandler) Route {
	return Route{router: router, id: uuid.New().String(), topic: topic, handler: handler}
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

func (r *Route) match(message *Message) bool {
	return r.topic == message.Topic() || routeIncludesTopic(r.topic, message.Topic())
}

func (r *Route) vars(message *Message) []string {
	var vars []string
	route := routeSplit(r.topic)
	topic := strings.Split(message.Topic(), "/")

	for i, section := range route {
		if section == "+" {
			if len(topic) > i {
				vars = append(vars, topic[i])
			}
		} else if section == "#" {
			if len(topic) > i {
				vars = append(vars, topic[i:]...)
			}
		}
	}

	return vars
}

func (r *router) addRoute(topic string, handler MessageHandler) Route {
	if handler != nil {
		route := newRoute(r, topic, handler)
		r.lock.Lock()
		r.routes = append(r.routes, route)
		r.lock.Unlock()
		return route
	}
	return Route{router: r}
}

func (r *router) removeRoute(removeRoute *Route) {
	r.lock.Lock()
	for i, route := range r.routes {
		if route.id == removeRoute.id {
			r.routes[i] = r.routes[len(r.routes)-1]
			r.routes = r.routes[:len(r.routes)-1]
		}
	}
	r.lock.Unlock()
}

func (r *router) match(message *Message) []Route {
	routes := []Route{}
	r.lock.RLock()
	for _, route := range r.routes {
		if route.match(message) {
			routes = append(routes, route)
		}
	}
	r.lock.RUnlock()
	return routes
}

// Stop removes this route from the router and stops matching it
func (r *Route) Stop() {
	r.router.removeRoute(r)
}
