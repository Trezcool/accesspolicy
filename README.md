# access-policy

This initiative employs a declarative and explicit strategy for handling access control in Go projects. 
It is consolidated in a single location and presented in a manner that is comprehensible to individuals with lesser technical expertise. 
If you have experience with other declarative access frameworks, like AWS' IAM, you will find the syntax to be familiar.

Example:

```go
package main

const rootUserID = 1

func isRoot(user User, action Action) bool {
	return user.GetID() == rootUserID
}

func main() {
	// Define a policy
	policy := AccessPolicy{
		Statements: []Statement{
			{
				Actions:   Actions{ActionAll},
				Principal: PrincipalAuthenticated,
				Conditions: Conditions{
					isRoot,
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
	usr := &user{id: rootUserID}
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
func (u *user) IsAnonymous() bool        { return u.id == 0 }
```
