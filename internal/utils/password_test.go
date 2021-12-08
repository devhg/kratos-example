package utils_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/devhg/kratos-example/internal/utils"
)

func TestGenerate(t *testing.T) {
	testPassword := "password"
	s := utils.GeneratePassword(testPassword)

	// 生成的密码与预期的一致
	assert.Equal(t, "c52f65639a16da778bd8839424495012", s)

	// 生成两次的密码保持一致
	assert.Equal(t, utils.GeneratePassword(testPassword), utils.GeneratePassword(testPassword))

	assert.Len(t, utils.GeneratePassword(testPassword), 32, "密码必须是32位")
}
