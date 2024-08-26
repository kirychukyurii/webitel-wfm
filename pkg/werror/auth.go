package werror

import (
	"fmt"
)

type AuthInvalidTokenError struct {
	id     string
	Reason error
}

func NewAuthInvalidTokenError(id string, reason error) AuthInvalidTokenError {
	return AuthInvalidTokenError{
		id:     id,
		Reason: reason,
	}
}

func (e AuthInvalidTokenError) Id() string {
	return e.id
}

func (e AuthInvalidTokenError) RPCError() string {
	return fmt.Sprintf("auth token is invalid: %s", e.Reason)
}

func (e AuthInvalidTokenError) Error() string {
	return fmt.Sprintf("%s: %s", e.Id(), e.RPCError())
}

type AuthInvalidSessionError struct {
	id     string
	Token  string
	Reason error
}

func NewAuthInvalidSessionError(id, token string, reason error) AuthInvalidSessionError {
	return AuthInvalidSessionError{
		id:     id,
		Token:  token,
		Reason: reason,
	}
}

func (e AuthInvalidSessionError) Id() string {
	return e.id
}

func (e AuthInvalidSessionError) RPCError() string {
	return fmt.Sprintf("auth session [%s] is invalid: %s", e.Token, e.Reason)

}

func (e AuthInvalidSessionError) Error() string {
	return fmt.Sprintf("%s: %s", e.Id(), e.RPCError())
}

type AuthForbiddenError struct {
	id     string
	Scope  string
	Action string
}

func NewAuthForbiddenError(id, scope, action string) AuthForbiddenError {
	return AuthForbiddenError{
		id:     id,
		Scope:  scope,
		Action: action,
	}
}

func (e AuthForbiddenError) Id() string {
	return e.id
}

func (e AuthForbiddenError) RPCError() string {
	return fmt.Sprintf("permission %s denied on resource %s (or it might not exist)", e.Action, e.Scope)
}

func (e AuthForbiddenError) Error() string {
	return fmt.Sprintf("%s: %s", e.Id(), e.RPCError())
}
