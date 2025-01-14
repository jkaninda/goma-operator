package controller

// RouteByName implements sort.Interface for []Route based on the Name field.
type RouteByName []Route

// MiddlewareByName implements sort.Interface for []Middleware based on the Name field.
type MiddlewareByName []Middleware

func (a RouteByName) Len() int           { return len(a) }
func (a RouteByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RouteByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (a MiddlewareByName) Len() int           { return len(a) }
func (a MiddlewareByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MiddlewareByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
