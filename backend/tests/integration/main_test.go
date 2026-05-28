package integration

import (
	"main/internal/testutil"
	"os"
	"testing"
)

var suite *testutil.TestSuite

func TestMain(m *testing.M) {
	suite = testutil.SetupSuiteForTestMain()
	code := m.Run()
	suite.Teardown()
	os.Exit(code)
}
