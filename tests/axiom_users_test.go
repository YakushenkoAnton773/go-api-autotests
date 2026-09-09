package tests

import (
	"testing"

	"github.com/Nikita-Filonov/axiom"
	"github.com/stretchr/testify/assert"
)

type axiomUser struct {
	Email  string
	Active bool
}

func TestAxiomCreateUser(t *testing.T) {
	testRunner := axiom.NewRunner()
	testCase := axiom.NewCase(
		axiom.WithCaseName("user can be created"),
	)

	testRunner.RunCase(t, testCase, func(cfg *axiom.Config) {
		user := axiomUser{
			Email:  "student@example.com",
			Active: true,
		}

		assert.Equal(cfg.T(), "student@example.com", user.Email)
		assert.True(cfg.T(), user.Active)
	})
}

func TestAxiomDeactivateUser(t *testing.T) {
	testRunner := axiom.NewRunner()
	testCase := axiom.NewCase(
		axiom.WithCaseName("user can be deactivated"),
	)

	testRunner.RunCase(t, testCase, func(cfg *axiom.Config) {
		user := axiomUser{
			Email:  "student@example.com",
			Active: true,
		}

		user.Active = false

		assert.False(cfg.T(), user.Active)
	})
}
