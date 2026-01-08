package entra

// User represents an Entra ID (Azure AD) user
type User struct {
	ID                string   `json:"id" yaml:"id"`
	DisplayName       string   `json:"displayName" yaml:"display_name"`
	UserPrincipalName string   `json:"userPrincipalName" yaml:"user_principal_name"`
	Mail              string   `json:"mail,omitempty" yaml:"mail,omitempty"`
	JobTitle          string   `json:"jobTitle,omitempty" yaml:"job_title,omitempty"`
	Department        string   `json:"department,omitempty" yaml:"department,omitempty"`
	Roles             []string `json:"roles,omitempty" yaml:"roles,omitempty"`
	Groups            []string `json:"groups,omitempty" yaml:"groups,omitempty"`
}

// ServicePrincipal represents an Entra ID service principal
type ServicePrincipal struct {
	ID          string   `json:"id" yaml:"id"`
	AppID       string   `json:"appId" yaml:"app_id"`
	DisplayName string   `json:"displayName" yaml:"display_name"`
	Secrets     []Secret `json:"-" yaml:"secrets,omitempty"` // Hidden from JSON output
	Permissions []string `json:"permissions,omitempty" yaml:"permissions,omitempty"`
}

// Secret represents a service principal secret/credential
type Secret struct {
	KeyID       string `json:"keyId" yaml:"key_id"`
	DisplayName string `json:"displayName" yaml:"display_name"`
	EndDateTime string `json:"endDateTime" yaml:"end_date_time"`
	SecretText  string `json:"-" yaml:"secret_text"` // The actual secret value
}

// Group represents an Entra ID group
type Group struct {
	ID          string   `json:"id" yaml:"id"`
	DisplayName string   `json:"displayName" yaml:"display_name"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Members     []string `json:"members,omitempty" yaml:"members,omitempty"`
}

// Environment holds the mocked Entra ID state for a scenario
type Environment struct {
	Users             []User
	ServicePrincipals []ServicePrincipal
	Groups            []Group
}

// NewEnvironment creates a new Entra ID environment from scenario resources
func NewEnvironment(users []User, servicePrincipals []ServicePrincipal, groups []Group) *Environment {
	return &Environment{
		Users:             users,
		ServicePrincipals: servicePrincipals,
		Groups:            groups,
	}
}

// ListUsers returns all users
func (e *Environment) ListUsers() []User {
	if e == nil || e.Users == nil {
		return []User{}
	}
	return e.Users
}

// GetUser returns a user by UPN or ID
func (e *Environment) GetUser(identifier string) *User {
	if e == nil {
		return nil
	}
	for i, u := range e.Users {
		if u.UserPrincipalName == identifier || u.ID == identifier {
			return &e.Users[i]
		}
	}
	return nil
}

// ListServicePrincipals returns all service principals
func (e *Environment) ListServicePrincipals() []ServicePrincipal {
	if e == nil || e.ServicePrincipals == nil {
		return []ServicePrincipal{}
	}
	return e.ServicePrincipals
}

// GetServicePrincipal returns a service principal by app ID or object ID
func (e *Environment) GetServicePrincipal(identifier string) *ServicePrincipal {
	if e == nil {
		return nil
	}
	for i, sp := range e.ServicePrincipals {
		if sp.AppID == identifier || sp.ID == identifier {
			return &e.ServicePrincipals[i]
		}
	}
	return nil
}

// ListGroups returns all groups
func (e *Environment) ListGroups() []Group {
	if e == nil || e.Groups == nil {
		return []Group{}
	}
	return e.Groups
}

// GetGroup returns a group by ID or display name
func (e *Environment) GetGroup(identifier string) *Group {
	if e == nil {
		return nil
	}
	for i, g := range e.Groups {
		if g.ID == identifier || g.DisplayName == identifier {
			return &e.Groups[i]
		}
	}
	return nil
}
