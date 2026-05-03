package validate_test

import (
	"testing"
	"xcross-go-starter/utility/validate"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestIsEmail(t *testing.T) {
	b := validate.IsEmail("QTT123456@163.com")
	gtest.Assert(true, b)
}
