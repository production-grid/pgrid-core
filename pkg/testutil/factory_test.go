package testutil

import (
	"testing"

	"github.com/production-grid/pgrid-core/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestBasisPointGeneration(t *testing.T) {

	bps := GenerateTestBasisPoints(40)

	assert := assert.New(t)
	assert.NotNil(bps)

	logging.Infoln(bps.String())

}
