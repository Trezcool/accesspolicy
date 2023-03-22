# access-policy

This initiative employs a declarative and explicit strategy for handling access control in Go projects. 
It is consolidated in a single location and presented in a manner that is comprehensible to individuals with lesser technical expertise. 
If you have experience with other declarative access frameworks, like AWS' IAM, you will find the syntax to be familiar.

Example:

```go
package main

func main() {
	// Define a policy
	policy := AccessPolicy{
		Statements: []Statement{
			{
				Actions:   Actions{ActionAll},
				Principal: PrincipalAuthenticated,
				Conditions: Conditions{
					func(user User, action Action) bool {
						return user.GetID() == 1
					},
				},
				Effect: EffectAllow,
			},
			{
				Actions:   Actions{ActionAnySafe},
				Principal: PrincipalAuthenticated,
				Effect:    EffectAllow,
			},
		},
	}

	// Define a user and an action
	usr := &user{id: 1}
	action := HTTPMethodAction(http.MethodGet)

	// Enforce the policy
	if policy.HasPermission(usr, action) {
		// Allow
	} else {
		// Deny
	}
}

type user struct{ id uint }
func (u *user) GetID() uint              { return u.id }
func (u *user) GetGroups() []string      { return nil }
func (u *user) GetPermissions() []string { return nil }
func (u *user) IsAnonymous() bool        { return false }
```
