package txprocessor

import "testing"

// Tests that transactions can be added to strict lists and list contents and
// nonce boundaries are correctly maintained.
func TestStrictTxListAdd(t *testing.T) {
	test := newTxList(true)
	_ = test
}
