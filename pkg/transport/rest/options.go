package rest

func WithBody[T any](body T) ResponseOpt[T] {
	return func(res *response[T]) {
		res.body = body
	}
}

func WithStatusCode[T any](statusCode int, _ T) ResponseOpt[T] {
	return func(res *response[T]) {
		res.statusCode = statusCode
	}
}
