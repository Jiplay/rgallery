package server_test

import (
	"testing"

	"github.com/robbymilo/rgallery/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestGetHash(t *testing.T) {
	var ans uint32 = 3884452138
	assert.EqualValues(t, server.GetHash("3884452138"), ans, "they should be equal")
}

func TestDecodeURL(t *testing.T) {
	decoded, _ := server.DecodeURL("20230506-kri%c5%a1ka-gora")
	assert.EqualValues(t, decoded, "20230506-kriška-gora", "they should be equal")
}
