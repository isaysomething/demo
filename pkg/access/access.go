package access

import "github.com/casbin/casbin/v2"

type AccessManager struct {
	enforcer *casbin.Enforcer
}

func NewAccessManager(enforcer *casbin.Enforcer) *AccessManager {
	return &AccessManager{
		enforcer: enforcer,
	}
}

func (am *AccessManager) Can() {

}

func (am *AccessManager) GetAllRoles() []string {
	return am.enforcer.GetAllRoles()
}
