// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate2FASecret(t *testing.T) {
	secret, err := Generate2FASecret("101645075095748608")
	assert.Nil(t, err)
	assert.Len(t, secret, 32)
}

func TestVerify2FA(t *testing.T) {
	_, err := Generate2FASecret("101645075095748608")
	assert.Nil(t, err)
	assert.False(t, Verify2FA("101645075095748608", "12345678"))
}
