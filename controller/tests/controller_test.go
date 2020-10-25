package tests

import (
	"advertisement_crud/controller"
	"advertisement_crud/etc/config"
	"advertisement_crud/mocks"
	"io"
	"io/ioutil"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// APITestSuite делает Setup и Teardown для тестов API.
type UnitTestSuite struct {
	suite.Suite

	// These fields recreates for each test in BeforeTest()
	controller *controller.Controller
	app        *mocks.IApplication
}

func (s *UnitTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	// minimal configuration which required for start controller

	s.app = &mocks.IApplication{}

	s.controller = controller.New(
		s.app,
		config.Main{
			FQDN: "localhost:8080",
		},
	)
	s.controller.InitRoutes()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

// ReadAllAndAssertErr reads the contents and assert no error.
func (s *UnitTestSuite) ReadAllAndAssertErr(r io.Reader) string {
	content, err := ioutil.ReadAll(r)
	s.NoError(err, "read: %s", string(content))
	return string(content)
}
