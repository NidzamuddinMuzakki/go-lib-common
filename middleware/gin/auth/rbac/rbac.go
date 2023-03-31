package rbac

import (
	"bitbucket.org/moladinTech/go-lib-rbac/action"
	"bitbucket.org/moladinTech/go-lib-rbac/model"
	"github.com/go-playground/validator/v10"
)

type IClientRBAC interface {
	IsRoleAllowed(allowedRoles map[string]bool, token string, applicationCode string) (
		isAllowed bool, userDetail *model.Me, err error)
	IsPermissionAllowed(allowedPermissions map[string]bool, token string, applicationCode string) (
		isAllowed bool, userDetail *model.Me, err error)
	GetUserDetail(token string, applicationCode string) (userDetail *model.Me, err error)
}

type ClientRBAC struct {
	httpHost     string `validate:"required"`
	grpcHost     string `validate:"required"`
	RBACInstance action.Action
}

type Option func(*ClientRBAC)

func WithHTTPHost(host string) Option {
	return func(client *ClientRBAC) {
		client.httpHost = host
	}
}

func WithGRPCHost(host string) Option {
	return func(client *ClientRBAC) {
		client.grpcHost = host
	}
}

func NewRBAC(validator *validator.Validate, options ...Option) *ClientRBAC {
	clientOptions := ClientRBAC{}
	for idx := 0; idx < len(options); idx++ {
		optionFunc := options[idx]
		optionFunc(&clientOptions)
	}

	err := validator.Struct(clientOptions)
	if err != nil {
		panic(err.Error())
	}

	clientOptions.RBACInstance = action.Init(clientOptions.httpHost, clientOptions.grpcHost)
	return &clientOptions
}

func (c *ClientRBAC) IsRoleAllowed(allowedRoles map[string]bool, token string, applicationCode string) (
	isAllowed bool, userDetail *model.Me, err error) {
	userDetail, err = c.RBACInstance.Me(token, applicationCode)
	if err != nil {
		return false, nil, err
	}

	for idx := 0; idx < len(userDetail.Roles); idx++ {
		if _, isAllowed := allowedRoles[userDetail.Roles[idx]]; isAllowed {
			return isAllowed, userDetail, nil
		}
	}

	return false, userDetail, nil
}

func (c *ClientRBAC) IsPermissionAllowed(allowedPermissions map[string]bool, token string, applicationCode string) (
	isAllowed bool, userDetail *model.Me, err error) {
	userDetail, err = c.RBACInstance.Me(token, applicationCode)
	if err != nil {
		return false, nil, err
	}

	for idx := 0; idx < len(userDetail.Permissions); idx++ {
		if _, isAllowed := allowedPermissions[userDetail.Permissions[idx]]; isAllowed {
			return isAllowed, userDetail, nil
		}
	}

	return false, userDetail, nil
}

func (c *ClientRBAC) GetUserDetail(token string, applicationCode string) (userDetail *model.Me, err error) {
	userDetail, err = c.RBACInstance.Me(token, applicationCode)
	if err != nil {
		return nil, err
	}

	return userDetail, nil
}
