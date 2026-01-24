package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type environment struct {
	values map[string]any
}

func NewBaseEnviroment() *environment {
	return &environment{
		values: make(map[string]any),
	}
}

func (e *environment) define(name string, value any) {
	e.values[name] = value
}

func (e *environment) get(name *token.Token) (any, *RuntimeError) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}
	return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable %s.", name.Lexeme))
}
