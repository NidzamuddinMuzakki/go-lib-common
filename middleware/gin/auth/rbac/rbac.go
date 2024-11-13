package rbac

import (
	"bitbucket.org/moladinTech/go-lib-rbac/action"
	"bitbucket.org/moladinTech/go-lib-rbac/model"
	"github.com/go-playground/validator/v10"
)

type IClientRBAC interface {
	IsRoleAllowed(allowedRoles []string, token string) (isAllowed bool, userDetail *model.Me, err error)
	IsPermissionAllowed(allowedPermissions []string, token string) (isAllowed bool, userDetail *model.Me, err error)
	GetUserDetail(token string) (userDetail *model.Me, err error)
}

type ClientRBAC struct {
	httpHost        string `validate:"required"`
	grpcHost        string `validate:"required"`
	action          action.Action
	ApplicationCode string `validate:"required"`
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

func WithApplicationCode(applicationCode string) Option {
	return func(client *ClientRBAC) {
		client.ApplicationCode = applicationCode
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

	clientOptions.action = action.Init(clientOptions.httpHost, clientOptions.grpcHost)
	return &clientOptions
}

func (c *ClientRBAC) IsRoleAllowed(allowedRolesArr []string, token string) (isAllowed bool, userDetail *model.Me, err error) {
	allowedRoles := c.convertRolesToMap(allowedRolesArr...)
	userDetail, err = c.action.Me(token, c.ApplicationCode)
	if err != nil {
		return false, nil, err
	}

	for idx := 0; idx < len(userDetail.Roles); idx++ {
		if _, ok := allowedRoles[userDetail.Roles[idx]]; ok {
			return ok, userDetail, nil
		}
	}

	return false, userDetail, nil
}

func (c *ClientRBAC) convertRolesToMap(roles ...string) map[string]bool {
	allowedRoles := make(map[string]bool, len(roles))
	for _, role := range roles {
		allowedRoles[role] = true
	}

	return allowedRoles
}

func (c *ClientRBAC) IsPermissionAllowed(allowedPermissionsArr []string, token string) (isAllowed bool, userDetail *model.Me, err error) {
	allowedPermissions := c.convertPermissionToMap(allowedPermissionsArr...)
	userDetail, err = c.action.Me(token, c.ApplicationCode)
	if err != nil {
		return false, nil, err
	}

	for idx := 0; idx < len(userDetail.Permissions); idx++ {
		if _, ok := allowedPermissions[userDetail.Permissions[idx]]; ok {
			return ok, userDetail, nil
		}
	}

	return false, userDetail, nil
}

func (c *ClientRBAC) convertPermissionToMap(permissions ...string) map[string]bool {
	allowedPermissions := make(map[string]bool, len(permissions))
	for _, role := range permissions {
		allowedPermissions[role] = true
	}

	return allowedPermissions
}

func (c *ClientRBAC) GetUserDetail(token string) (userDetail *model.Me, err error) {
	userDetail, err = c.action.Me(token, c.ApplicationCode)
	if err != nil {
		return nil, err
	}

	return userDetail, nil
}
