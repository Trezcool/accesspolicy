# accesspolicy

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]
[![version-img]][version-url]
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FTrezcool%2Faccesspolicy.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FTrezcool%2Faccesspolicy?ref=badge_shield)

This initiative employs a declarative and explicit strategy for handling access control in Go projects. 
It is consolidated in a single location and presented in a manner that is comprehensible to individuals with lesser technical expertise. 
If you have experience with other declarative access frameworks, like AWS' IAM, you will find the syntax to be familiar.

Example:

```go
package main

const rootUserID = 1

func isRoot(ctx context.Context, user User, action Action) bool {
	return user.GetID() == rootUserID
}

func main() {
	// Define a policy
	policy := Policy{
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
	ctx := context.Background()
	usr := &user{id: rootUserID}
	action := HTTPMethodAction(http.MethodGet)

	// Enforce the policy
	if policy.HasPermission(ctx, usr, action) {
		// Allow
	} else {
		// Deny
	}
}

type user struct{ id uint }
func (u *user) GetID() uint              { return u.id }
func (u *user) IsAnonymous() bool        { return u.id == 0 }
```

[build-img]: https://github.com/Trezcool/accesspolicy/workflows/ci/badge.svg
[build-url]: https://github.com/Trezcool/accesspolicy/actions
[pkg-img]: https://pkg.go.dev/badge/Trezcool/accesspolicy/v0
[pkg-url]: https://pkg.go.dev/github.com/Trezcool/accesspolicy/v0
[reportcard-img]: https://goreportcard.com/badge/Trezcool/accesspolicy
[reportcard-url]: https://goreportcard.com/report/Trezcool/accesspolicy
[coverage-img]: https://codecov.io/gh/Trezcool/accesspolicy/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/Trezcool/accesspolicy
[version-img]: https://img.shields.io/github/v/release/Trezcool/accesspolicy
[version-url]: https://github.com/Trezcool/accesspolicy/releases


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FTrezcool%2Faccesspolicy.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FTrezcool%2Faccesspolicy?ref=badge_large)