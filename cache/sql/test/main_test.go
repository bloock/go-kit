package test

import (
	"github.com/bloock/go-kit/test_utils"
	"testing"
)

func TestMain(m *testing.M) {
	test_utils.SetupPostgresIntegrationTest(m, 120, "../migrations")
}
