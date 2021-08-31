package sqlcommenter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentEscapse(t *testing.T) {
	tags := Tags{
		"route":     `/param first`,
		"num":       `1234`,
		"query":     `DROP TABLE FOO'`,
		"injection": `/route/*/;DROP TABLE USERS`,
	}

	assert.Equal(t, `injection='%2Froute%2F%2A%2F%3BDROP%20TABLE%20USERS',num='1234',query='DROP%20TABLE%20FOO%27',route='%2Fparam%20first'`, tags.Marshal())
}
