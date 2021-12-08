// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package utils_test

import (
	"github.com/devhg/kratos-example/internal/utils"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRandomString(t *testing.T) {
	assert.Len(t, utils.RandomString(1), 1)
	assert.Len(t, utils.RandomString(2), 2)
	assert.Len(t, utils.RandomString(3), 3)
	assert.Len(t, utils.RandomString(4), 4)
	assert.Len(t, utils.RandomString(8), 8)
	assert.Len(t, utils.RandomString(16), 16)
	assert.IsType(t, "string", utils.RandomString(16))
}

func TestRandomNumeric(t *testing.T) {
	assert.Len(t, utils.RandomNumeric(1), 1)
	assert.Len(t, utils.RandomNumeric(2), 2)
	assert.Len(t, utils.RandomNumeric(3), 3)
	assert.Len(t, utils.RandomNumeric(4), 4)
	assert.Len(t, utils.RandomNumeric(8), 8)
	assert.Len(t, utils.RandomNumeric(16), 16)
	assert.IsType(t, "string", utils.RandomNumeric(16))
	assert.True(t, regexp.MustCompile(`^\d+$`).MatchString(utils.RandomNumeric(32)))
}
