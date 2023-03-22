package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type (
	User interface {
		IsAnonymous() bool
	}
	userWithID interface {
		User
		GetID() uint
	}
	userWithGroups interface {
		User
		GetGroups() []string
	}
	userWithPermissions interface {
		User
		GetPermissions() []string
	}
)

type AccessPolicy struct {
	Statements []Statement
}

type Statement struct {
	Actions    Actions
	Principal  Principal
	Conditions Conditions
	Effect     Effect
}

func (p *AccessPolicy) HasPermission(user User, action Action) bool {
	if len(p.Statements) == 0 {
		return false
	}
	return p.evaluateStatements(user, action)
}

func (p *AccessPolicy) evaluateStatements(user User, action Action) bool {
	matched := p.getStatementsMatchingAction(action)
	matched = p.getStatementsMatchingPrincipal(matched, user)
	matched = p.getStatementsMatchingConditions(matched, user, action)

	denied := p.getDeniedStatements(matched)

	if len(matched) == 0 || len(denied) > 0 {
		return false
	}
	return true
}

func (p *AccessPolicy) getStatementsMatchingAction(action Action) []Statement {
	return lo.Filter(p.Statements, func(statement Statement, _ int) bool {
		return statement.Actions.Match(action)
	})
}
func (p *AccessPolicy) getStatementsMatchingPrincipal(statements []Statement, user User) []Statement {
	return lo.Filter(statements, func(statement Statement, _ int) bool {
		return statement.Principal.Match(user)
	})
}
func (p *AccessPolicy) getStatementsMatchingConditions(statements []Statement, user User, action Action) []Statement {
	return lo.Filter(statements, func(statement Statement, _ int) bool {
		return statement.Conditions.Match(user, action)
	})
}
func (p *AccessPolicy) getDeniedStatements(statements []Statement) []Statement {
	return lo.Filter(statements, func(statement Statement, _ int) bool {
		return statement.Effect != EffectAllow
	})
}

type Action struct {
	Name   string
	IsSafe bool
}

type Actions []Action

func (l Actions) Match(action Action) bool {
	switch {
	case lo.Contains(l, ActionAll):
		return true
	case lo.Contains(l, ActionAnySafe) && action.IsSafe:
		return true
	}
	return lo.Contains(l, action)
}

var (
	ActionAll     = Action{"*", false}
	ActionAnySafe = Action{"any_safe", true}
)

var safeHTTPMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
}

func HTTPMethodAction(method string) Action {
	return Action{method, lo.Contains(safeHTTPMethods, method)}
}

type Principal string

func (p Principal) Match(user User) bool {
	switch {
	case p == PrincipalAll:
		return true
	case p == PrincipalAuthenticated:
		return !user.IsAnonymous()
	case p == PrincipalAnonymous:
		return user.IsAnonymous()
	case strings.HasPrefix(string(p), principalGroupPrefix):
		return p.matchGroups(user)
	case strings.HasPrefix(string(p), principalPermissionPrefix):
		return p.matchPermissions(user)
	case strings.HasPrefix(string(p), principalUserPrefix):
		return p.matchUser(user)
	default:
		return false
	}
}

func (p Principal) matchGroups(user User) bool {
	u, ok := user.(userWithGroups)
	if !ok {
		return false
	}
	pGroupsStr := strings.TrimPrefix(string(p), principalGroupPrefix)
	pGroups := strings.Split(pGroupsStr, ",")
	uGroups := u.GetGroups()
	return len(lo.Intersect(pGroups, uGroups)) > 0
}
func (p Principal) matchPermissions(user User) bool {
	u, ok := user.(userWithPermissions)
	if !ok {
		return false
	}
	pPermsStr := strings.TrimPrefix(string(p), principalPermissionPrefix)
	pPerms := strings.Split(pPermsStr, ",")
	uPerms := u.GetPermissions()
	return len(lo.Intersect(pPerms, uPerms)) == len(pPerms)
}
func (p Principal) matchUser(user User) bool {
	u, ok := user.(userWithID)
	if !ok {
		return false
	}
	pUsersStr := strings.TrimPrefix(string(p), principalUserPrefix)
	pUsers := strings.Split(pUsersStr, ",")
	uID := strconv.Itoa(int(u.GetID()))
	return lo.Contains(pUsers, uID)
}

const (
	PrincipalAll           Principal = "*"
	PrincipalAuthenticated           = "authenticated"
	PrincipalAnonymous               = "anonymous"
)

const (
	principalGroupPrefix      = "group:"
	principalPermissionPrefix = "permission:"
	principalUserPrefix       = "id:"
)

// GroupPrincipal will match any user that is in any of the groups
func GroupPrincipal(group ...string) Principal {
	return Principal(principalGroupPrefix + strings.Join(group, ","))
}

// PermissionPrincipal will match any user that has all the permissions TODO: support OR ?
func PermissionPrincipal(permission ...string) Principal {
	return Principal(principalPermissionPrefix + strings.Join(permission, ","))
}

// UserPrincipal will match any user whose ID is in the list
func UserPrincipal(user ...string) Principal {
	return Principal(principalUserPrefix + strings.Join(user, ","))
}

type Condition func(user User, action Action) bool

type Conditions []Condition

func (l Conditions) Match(user User, action Action) bool {
	return lo.EveryBy(l, func(condition Condition) bool {
		return condition(user, action)
	})
}

type Effect string

const (
	EffectAllow Effect = "allow"
	EffectDeny         = "deny"
)
