package migrate

import (
	"testing"
)

func TestCommitOrRollback(t *testing.T) {
	Connect()
	DeleteRowByTable("1", "2", "")
}
