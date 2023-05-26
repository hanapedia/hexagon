package constants

// GetHttpMethodFromAction converts Action to HttpMethod
func GetHttpMethodFromAction(action Action) HttpMethod {
	switch action {
	case "read":
		return GET
	case "write":
		return POST
	default:
		return GET
	}
}
