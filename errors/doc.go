// Package errors contains error definitions and utilities
// Application Error: This is the base error. If you have a more specific error which fits the scenario, you should use it.
// InternalServer Error: Use when you have a logical error on the server
// BadRequest Error: Use it when responsing to a transport and the input data do not match api expectations
// Forbidden Error: (403 http status) Use it when the user is authenticated by the request is forbidden for him
// Unauthotized Error: (401 http status) Use it when the user is not authenticated thus not authorized to perform the operation
// MethodNotAllowed Error: (405 http status) Use it when method was received and recognized by the server, but the server has rejected that particular method for the requested resource.
// Timeout Error: (408) Use it to indicate that the time allowed to complete the operation completed. You can use it for network calls, API request, database calls and etc.
// NotFound Error: (404) Use it to indicates that the server itself was found, but that the server was not able to retrieve the requested data. For example, the user was not found.
// Database Error: Use when you have a database related error
package errors
