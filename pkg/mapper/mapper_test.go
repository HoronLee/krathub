package mapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type domainUser struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Phone    *string
	Role     string
}

type entLikeUser struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Phone    *string
	Role     string
}

func TestCopierMapperWithEntLikeStruct(t *testing.T) {
	phone := "13800000000"
	domain := &domainUser{
		ID:       7,
		Name:     "alice",
		Email:    "alice@example.com",
		Password: "hashed-password",
		Phone:    &phone,
		Role:     "admin",
	}

	m := New[domainUser, entLikeUser]().RegisterConverters(AllBuiltinConverters())

	entity := m.ToEntity(domain)
	require.NotNil(t, entity)
	require.Equal(t, domain.ID, entity.ID)
	require.Equal(t, domain.Name, entity.Name)
	require.Equal(t, domain.Email, entity.Email)
	require.Equal(t, domain.Password, entity.Password)
	require.Equal(t, domain.Phone, entity.Phone)
	require.Equal(t, domain.Role, entity.Role)

	back := m.ToDomain(entity)
	require.NotNil(t, back)
	require.Equal(t, domain, back)
}

func TestCopierMapperListWithNilItems(t *testing.T) {
	phone := "13900000000"
	entities := []*entLikeUser{
		{ID: 1, Name: "u1", Email: "u1@example.com", Phone: &phone, Role: "user"},
		nil,
		{ID: 2, Name: "u2", Email: "u2@example.com", Role: "admin"},
	}

	m := New[domainUser, entLikeUser]().RegisterConverters(AllBuiltinConverters())
	domains := m.ToDomainList(entities)
	require.Len(t, domains, 2)
	require.Equal(t, int64(1), domains[0].ID)
	require.Equal(t, int64(2), domains[1].ID)
}
