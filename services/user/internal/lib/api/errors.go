package api

func UnexpectedError() Response {
	return Err("an unexpected error occured.")
}
