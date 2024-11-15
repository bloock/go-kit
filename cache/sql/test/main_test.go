package test

import (
	"github.com/bloock/go-kit/test_utils/postgres"
	"testing"
)

func TestMain(m *testing.M) {
	postgres.SetupPostgresIntegrationTest(m, 120, "../migrations")
}
