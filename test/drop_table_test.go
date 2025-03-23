package gomb_test

import (
	"testing"

	gomb "github.com/nandrechetan/gomb/internal"
	"github.com/stretchr/testify/assert"
)

func TestDropTable(t *testing.T) {
	t.Run("Basic Index Creation", func(t *testing.T) {
		table := gomb.NewDropTable("account")
		sql, err := table.ToSQL()
		if err != nil {
			t.Error(err)
		}
		expectedSQL := "DROP TABLE IF EXISTS account"
		assert.Equal(t, expectedSQL, sql)
	})
}
