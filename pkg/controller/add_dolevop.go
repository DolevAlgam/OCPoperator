package controller

import (
	"github.com/opdolev/pkg/controller/dolevop"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, dolevop.Add)
}