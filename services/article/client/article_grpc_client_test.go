package client

import (
	"testing"

	"github.com/loak155/techbranch-backend/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewArticleGRPCClient(t *testing.T) {
	conf, err := config.Load("../../../configs/config.yaml")
	assert.NoError(t, err)
	client, err := NewArticleGRPCClient(conf)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
