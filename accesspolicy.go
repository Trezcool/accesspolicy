package accesspolicy

import (
	"context"
	"net/http"
	"strings"

	"github.com/samber/lo"
)

// todo: tests
// todo: ci
// todo: cd -> godoc with examples

type (
	User interface {
		IsAnonymous() bool
	}
	userWithID interface {
		User
		GetIDStr() string
	}
	userWithGroups interface {
		User
		GetGroups() []string
	}
	userWithPermissions interface {
		User
		GetPermissions() []string
	}
	superUser interface {
		User
		IsSuperUser() bool
	}
)

// Policy is a collection of statements that provide access to resources.
//
// A policy is evaluated by checking if the user has at least one statement
// that matches the action and effect is set to allow.
type Policy struct {
	Statements []Statement
}

// Statement is a single statement in a policy that provides access to resources.
type Statement struct {
	Actions    Actions
	Principal  Principal
	Conditions Conditions
	Effect     Effect
}

// HasPermission checks if the user has permission to perform the action.
func (p *Policy) HasPermission(ctx context.Context, user User, action Action) bool {
	if su, ok := user.(superUser); ok && su.IsSuperUser() {
		return true
	}

	if len(p.Statements) == 0 {
		return false
	}
	return p.evaluateStatements(ctx, user, action)
}

func (p *Policy) evaluateStatements(ctx context.Context, user User, action Action) bool {
	matched := p.getStatementsMatchingAction(action)
	matched = p.getStatementsMatchingPrincipal(matched, user)
	matched = p.getStatementsMatchingConditions(ctx, matched, user, action)

	denied := p.getDeniedStatements(matched)

	if len(matched) == 0 || len(denied) > 0 {
		return false
	}
	return true
}

func (p *Policy) getStatementsMatchingAction(action Action) []Statement {
	return lo.Filter(p.Statements, func(statement Statement, _ int) bool {
		return statement.Actions.Match(action)
	})
}
func (p *Policy) getStatementsMatchingPrincipal(statements []Statement, user User) []Statement {
	return lo.Filter(statements, func(statement Statement, _ int) bool {
		return statement.Principal.Match(user)
	})
}
func (p *Policy) getStatementsMatchingConditions(ctx context.Context, statements []Statement, user User, action Action) []Statement {
	return lo.Filter(statements, func(statement Statement, _ int) bool {
		return statement.Conditions.Match(ctx, user, action)
	})
}
func (p *Policy) getDeniedStatements(statements []Statement) []Statement {
	return lo.Filter(statements, func(statement Statement, _ int) bool {
		return statement.Effect != EffectAllow
	})
}

// Action represents an action that can be performed on a resource.
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

// HTTPMethodAction returns an Action for the given HTTP method.
func HTTPMethodAction(method string) Action {
	return Action{method, lo.Contains(safeHTTPMethods, method)}
}

// Principal represents a user or group that can perform an action.
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
	return len(lo.Intersect(pGroups, uGroups)) > 0 // ANY
}
func (p Principal) matchPermissions(user User) bool {
	u, ok := user.(userWithPermissions)
	if !ok {
		return false
	}
	pPermsStr := strings.TrimPrefix(string(p), principalPermissionPrefix)
	pPerms := strings.Split(pPermsStr, ",")
	uPerms := u.GetPermissions()
	return len(lo.Intersect(pPerms, uPerms)) == len(pPerms) // ALL
}
func (p Principal) matchUser(user User) bool {
	u, ok := user.(userWithID)
	if !ok {
		return false
	}
	pUIDsStr := strings.TrimPrefix(string(p), principalUserPrefix)
	pUIDs := strings.Split(pUIDsStr, ",")
	uID := u.GetIDStr()
	return lo.Contains(pUIDs, uID)
}

const (
	PrincipalAll           Principal = "*"
	PrincipalAuthenticated Principal = "authenticated"
	PrincipalAnonymous     Principal = "anonymous"
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

// PermissionPrincipal will match any user that has all the permissions
func PermissionPrincipal(permission ...string) Principal {
	//  TODO: support OR ?
	return Principal(principalPermissionPrefix + strings.Join(permission, ","))
}

// UserPrincipal will match any user whose ID is in the list
func UserPrincipal(userID ...string) Principal {
	return Principal(principalUserPrefix + strings.Join(userID, ","))
}

// Condition represents a condition that must be met for an action to be allowed.
type Condition func(ctx context.Context, user User, action Action) bool

type Conditions []Condition

func (l Conditions) Match(ctx context.Context, user User, action Action) bool {
	return lo.EveryBy(l, func(condition Condition) bool {
		return condition(ctx, user, action)
	})
}

// Effect represents the effect of a statement.
type Effect string

const (
	EffectAllow Effect = "allow"
	EffectDeny  Effect = "deny"
)
