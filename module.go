package main

import (
	"net/http"
)

// Module is the interface that all Modules must implement.
type Module interface {
	Register(h Context)
}

// Context is the module's context passed into the Module.Register() call.
type Context interface {
	AddAuthHttp(path string, handler web.Handler)
	AddHttpHandler(path string, handler http.Handler)
	RegisterEventHandler(h bus.NodeHandler)
	ResetData(reset func())
}
