package entra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironment(t *testing.T) {
	users := []User{
		{ID: "user1", DisplayName: "User One"},
	}
	sps := []ServicePrincipal{
		{ID: "sp1", DisplayName: "Service Principal One"},
	}
	groups := []Group{
		{ID: "group1", DisplayName: "Group One"},
	}

	env := NewEnvironment(users, sps, groups)

	assert.NotNil(t, env)
	assert.Len(t, env.Users, 1)
	assert.Len(t, env.ServicePrincipals, 1)
	assert.Len(t, env.Groups, 1)
}

func TestNewEnvironmentWithNil(t *testing.T) {
	env := NewEnvironment(nil, nil, nil)

	assert.NotNil(t, env)
	assert.Nil(t, env.Users)
	assert.Nil(t, env.ServicePrincipals)
	assert.Nil(t, env.Groups)
}

func TestListUsers(t *testing.T) {
	users := []User{
		{ID: "1", DisplayName: "User1"},
		{ID: "2", DisplayName: "User2"},
	}
	env := NewEnvironment(users, nil, nil)

	result := env.ListUsers()

	assert.Len(t, result, 2)
	assert.Equal(t, "User1", result[0].DisplayName)
}

func TestListUsersEmpty(t *testing.T) {
	env := NewEnvironment(nil, nil, nil)

	result := env.ListUsers()

	assert.Empty(t, result)
}

func TestListUsersNilEnvironment(t *testing.T) {
	var env *Environment

	result := env.ListUsers()

	assert.Empty(t, result)
}

func TestGetUser(t *testing.T) {
	users := []User{
		{
			ID:                "user-123",
			DisplayName:       "John Doe",
			UserPrincipalName: "john@contoso.com",
		},
	}
	env := NewEnvironment(users, nil, nil)

	t.Run("by ID", func(t *testing.T) {
		user := env.GetUser("user-123")
		assert.NotNil(t, user)
		assert.Equal(t, "John Doe", user.DisplayName)
	})

	t.Run("by UPN", func(t *testing.T) {
		user := env.GetUser("john@contoso.com")
		assert.NotNil(t, user)
		assert.Equal(t, "user-123", user.ID)
	})

	t.Run("not found", func(t *testing.T) {
		user := env.GetUser("nonexistent")
		assert.Nil(t, user)
	})

	t.Run("nil environment", func(t *testing.T) {
		var nilEnv *Environment
		user := nilEnv.GetUser("user-123")
		assert.Nil(t, user)
	})
}

func TestListServicePrincipals(t *testing.T) {
	sps := []ServicePrincipal{
		{ID: "sp1", DisplayName: "App1"},
		{ID: "sp2", DisplayName: "App2"},
	}
	env := NewEnvironment(nil, sps, nil)

	result := env.ListServicePrincipals()

	assert.Len(t, result, 2)
}

func TestListServicePrincipalsEmpty(t *testing.T) {
	env := NewEnvironment(nil, nil, nil)

	result := env.ListServicePrincipals()

	assert.Empty(t, result)
}

func TestListServicePrincipalsNilEnvironment(t *testing.T) {
	var env *Environment

	result := env.ListServicePrincipals()

	assert.Empty(t, result)
}

func TestGetServicePrincipal(t *testing.T) {
	sps := []ServicePrincipal{
		{
			ID:          "sp-object-id",
			AppID:       "app-id-123",
			DisplayName: "My App",
		},
	}
	env := NewEnvironment(nil, sps, nil)

	t.Run("by object ID", func(t *testing.T) {
		sp := env.GetServicePrincipal("sp-object-id")
		assert.NotNil(t, sp)
		assert.Equal(t, "My App", sp.DisplayName)
	})

	t.Run("by app ID", func(t *testing.T) {
		sp := env.GetServicePrincipal("app-id-123")
		assert.NotNil(t, sp)
		assert.Equal(t, "sp-object-id", sp.ID)
	})

	t.Run("not found", func(t *testing.T) {
		sp := env.GetServicePrincipal("nonexistent")
		assert.Nil(t, sp)
	})

	t.Run("nil environment", func(t *testing.T) {
		var nilEnv *Environment
		sp := nilEnv.GetServicePrincipal("sp-object-id")
		assert.Nil(t, sp)
	})
}

func TestListGroups(t *testing.T) {
	groups := []Group{
		{ID: "g1", DisplayName: "Admins"},
		{ID: "g2", DisplayName: "Users"},
	}
	env := NewEnvironment(nil, nil, groups)

	result := env.ListGroups()

	assert.Len(t, result, 2)
}

func TestListGroupsEmpty(t *testing.T) {
	env := NewEnvironment(nil, nil, nil)

	result := env.ListGroups()

	assert.Empty(t, result)
}

func TestListGroupsNilEnvironment(t *testing.T) {
	var env *Environment

	result := env.ListGroups()

	assert.Empty(t, result)
}

func TestGetGroup(t *testing.T) {
	groups := []Group{
		{
			ID:          "group-id-123",
			DisplayName: "Security Admins",
			Description: "Security administrators",
		},
	}
	env := NewEnvironment(nil, nil, groups)

	t.Run("by ID", func(t *testing.T) {
		group := env.GetGroup("group-id-123")
		assert.NotNil(t, group)
		assert.Equal(t, "Security Admins", group.DisplayName)
	})

	t.Run("by display name", func(t *testing.T) {
		group := env.GetGroup("Security Admins")
		assert.NotNil(t, group)
		assert.Equal(t, "group-id-123", group.ID)
	})

	t.Run("not found", func(t *testing.T) {
		group := env.GetGroup("nonexistent")
		assert.Nil(t, group)
	})

	t.Run("nil environment", func(t *testing.T) {
		var nilEnv *Environment
		group := nilEnv.GetGroup("group-id-123")
		assert.Nil(t, group)
	})
}

func TestUserFields(t *testing.T) {
	user := User{
		ID:                "user-uuid",
		DisplayName:       "Jane Smith",
		UserPrincipalName: "jane@contoso.com",
		Mail:              "jane.smith@contoso.com",
		JobTitle:          "Engineer",
		Department:        "IT",
		Roles:             []string{"Reader", "Contributor"},
		Groups:            []string{"group1", "group2"},
	}

	assert.Equal(t, "user-uuid", user.ID)
	assert.Equal(t, "Jane Smith", user.DisplayName)
	assert.Equal(t, "jane@contoso.com", user.UserPrincipalName)
	assert.Equal(t, "jane.smith@contoso.com", user.Mail)
	assert.Equal(t, "Engineer", user.JobTitle)
	assert.Equal(t, "IT", user.Department)
	assert.Len(t, user.Roles, 2)
	assert.Len(t, user.Groups, 2)
}

func TestServicePrincipalFields(t *testing.T) {
	sp := ServicePrincipal{
		ID:          "sp-uuid",
		AppID:       "app-uuid",
		DisplayName: "Backup Service",
		Secrets: []Secret{
			{
				KeyID:       "key1",
				DisplayName: "Primary Key",
				EndDateTime: "2025-01-01",
				SecretText:  "secret-value",
			},
		},
		Permissions: []string{"Storage.Read", "KeyVault.Read"},
	}

	assert.Equal(t, "sp-uuid", sp.ID)
	assert.Equal(t, "app-uuid", sp.AppID)
	assert.Equal(t, "Backup Service", sp.DisplayName)
	assert.Len(t, sp.Secrets, 1)
	assert.Equal(t, "secret-value", sp.Secrets[0].SecretText)
	assert.Len(t, sp.Permissions, 2)
}

func TestGroupFields(t *testing.T) {
	group := Group{
		ID:          "group-uuid",
		DisplayName: "DevOps Team",
		Description: "DevOps engineers",
		Members:     []string{"user1", "user2", "user3"},
	}

	assert.Equal(t, "group-uuid", group.ID)
	assert.Equal(t, "DevOps Team", group.DisplayName)
	assert.Equal(t, "DevOps engineers", group.Description)
	assert.Len(t, group.Members, 3)
}
