# go-errors

[![CircleCI](https://circleci.com/gh/weathersource/go-errors.svg?style=shield)](https://circleci.com/gh/weathersource/go-errors)
[![GoDoc](https://img.shields.io/badge/godoc-ref-blue.svg)](https://godoc.org/github.com/weathersource/go-errors) |

Package errors is a robust and easy-to-use error package. It supports verbosity levels, maintains causal errors, implements stack tracing,
maps to common HTTP error codes, exposes temporary and timeout errors, implements JSON marshalling, and is GRPC compatible.

This package started [github.com/goware/errorx](https://github.com/goware/errorx) and grew from there.

Supported errors are as follows:

Error                    |  HTTP Code                  |  Description
-----------------------------------------------------------------------------
AbortedError             |  409 CONFLICT               |  operation was aborted
AlreadyExistsError       |  409 CONFLICT               |  attempt to create an entity failed because one already exists
CanceledError            |  499 CLIENT CLOSED REQUEST  |  operation was canceled
DataLossError            |  500 INTERNAL SERVER ERROR  |  unrecoverable data loss or corruption
DeadlineExceededError    |  504 GATEWAY TIMEOUT        |  operation expired before completion
FailedPreconditionError  |  400 BAD REQUEST            |  operation rejected because system is not in a state required for operation's execution
InternalError            |  500 INTERNAL SERVER ERROR  |  some invariants expected by underlying system has been broken
InvalidArgumentError     |  400 BAD REQUEST            |  client specified an invalid argument
NotFoundError            |  404 NOT FOUND              |  requested entity was not found
NotImplementedError      |  501 NOT IMPLEMENTED        |  operation is not implemented
OutOfRangeError          |  400 BAD REQUEST            |  operation was attempted past the valid range
NewPassthroughError      |  varies                     |  passes 400 errors directly, 500 errors are mapped to InternalError
PermissionDeniedError    |  403 FORBIDDEN              |  the caller does not have permission to execute the specified operation
ResourceExhaustedError   |  429 TOO MANY REQUESTS      |  some resource has been exhausted
UnauthenticatedError     |  401 UNAUTHORIZED           |  the request does not have valid authentication credentials
UnavailableError         |  503 SERVICE UNAVAILABLE    |  the service is currently unavailable
UnknownError             |  500 INTERNAL SERVER ERROR  |  unknown server error