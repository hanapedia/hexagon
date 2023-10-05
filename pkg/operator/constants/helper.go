package constants

// GetHttpMethodFromAction converts Action to HttpMethod
func GetHttpMethodFromAction(action Action) HttpMethod {
	switch action {
	case "read":
		return HTTP_GET
	case "write":
		return HTTP_POST
	default:
		return HTTP_GET
	}
}
