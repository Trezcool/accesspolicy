package accesspolicy

import (
	"context"
	"reflect"
	"testing"
)

func TestAccessPolicy_HasPermission(t *testing.T) {
	t.Skip("todo")

	type fields struct {
		Statements []Statement
	}
	type args struct {
		ctx    context.Context
		user   User
		action Action
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Policy{
				Statements: tt.fields.Statements,
			}
			if got := p.HasPermission(tt.args.ctx, tt.args.user, tt.args.action); got != tt.want {
				t.Errorf("HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPMethodAction(t *testing.T) {
	t.Skip("todo")

	type args struct {
		method string
	}
	tests := []struct {
		name string
		args args
		want Action
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HTTPMethodAction(tt.args.method); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTTPMethodAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupPrincipal(t *testing.T) {
	t.Skip("todo")

	type args struct {
		group []string
	}
	tests := []struct {
		name string
		args args
		want Principal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupPrincipal(tt.args.group...); got != tt.want {
				t.Errorf("GroupPrincipal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionPrincipal(t *testing.T) {
	t.Skip("todo")

	type args struct {
		permission []string
	}
	tests := []struct {
		name string
		args args
		want Principal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PermissionPrincipal(tt.args.permission...); got != tt.want {
				t.Errorf("PermissionPrincipal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserPrincipal(t *testing.T) {
	t.Skip("todo")

	type args struct {
		userID []string
	}
	tests := []struct {
		name string
		args args
		want Principal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserPrincipal(tt.args.userID...); got != tt.want {
				t.Errorf("UserPrincipal() = %v, want %v", got, tt.want)
			}
		})
	}
}
