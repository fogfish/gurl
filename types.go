package gurl

/*

IOCat defines the category for abstract I/O with a side-effects
*/
type IOCat struct {
	Fail       error
	HTTP       *IOCatHTTP
	LogLevel   int
	sideEffect Arrow
}

/*

Unsafe applies a side effect on the category
*/
func (cat *IOCat) Unsafe() *IOCat {
	return cat.sideEffect(cat)
}

/*

Config defines configuration for the IO category
*/
type Config func(*IOCat) *IOCat

/*

Arrow is a morphism applied to IO category. The library implements:

↣ gurl/http/send, which defines writer morphism that focuses inside and
reshapes HTTP protocol request. The writer morphism is used to declare HTTP
method, destination URL, request headers and payload.

↣ gurl/http/recv, which defines reader morphism that focuses into side-effect,
HTTP protocol response. The reader morphism is a pattern matcher, is used to
match HTTP response code, headers and response payload.
*/
type Arrow func(*IOCat) *IOCat
