package api

import "fmt"

type MethodType string

const (
	MethodGet   MethodType = "Get"
	MethodPub   MethodType = "Pub"
	MethodSet   MethodType = "Set"
	MethodDel   MethodType = "Del"
	MethodCall  MethodType = "Call"
	MethodReply MethodType = "Reply"
)

var methodTypes = map[MethodType]struct{}{
	MethodGet:   {},
	MethodPub:   {},
	MethodSet:   {},
	MethodDel:   {},
	MethodCall:  {},
	MethodReply: {},
}

func ParseMethodType(s string) (*MethodType, error) {
	mType := MethodType(s)
	_, ok := methodTypes[mType]
	if !ok {
		return nil, fmt.Errorf(`cannot parse:[%s] as MethodType`, s)
	}
	return &mType, nil
}
