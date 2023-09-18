package graphml

import "errors"

var (
	ErrEdgeDoesNotExist  = errors.New("edge does not exist")
	ErrNodeDoesNotExist  = errors.New("node does not exist")
	ErrEdgeAlreadyExists  = errors.New("edge already exists")
	ErrNodeDataAlreadySet = errors.New("node data already set, will not override")
	ErrEdgeDataAlreadySet = errors.New("edge data already set, will not override")
)
