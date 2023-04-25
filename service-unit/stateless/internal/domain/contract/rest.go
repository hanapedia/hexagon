package contract

type RestRequestBody struct {
	Message string  `json:"message"`
	Payload *string `json:"payload"`
}

type RestResponseBody struct {
	Message string  `json:"message"`
	Payload *string `json:"payload"`
}

