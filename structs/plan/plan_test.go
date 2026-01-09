package plan

import (
	"testing"

	"github.com/homedepot/trainer/structs/transaction"
)

func TestLoadTxnFile_RejectsPathTraversal(t *testing.T) {
	p := &Plan{}
	
	tests := []struct {
		name string
		file string
	}{
		{"traversal up", "../../../etc/passwd"},
		{"middle traversal", "transactions/../../secrets.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
err := p.LoadTxnFile(TxnInclude{File: tt.file})
if err == nil {
t.Error("LoadTxnFile() should reject path traversal")
}
})
	}
}

func TestFindTransaction(t *testing.T) {
	p := &Plan{
		Txn: []transaction.Transaction{{Name: "txn1"}, {Name: "txn2"}},
	}

	txn, err := p.FindTransaction("txn1")
	if err != nil || txn.Name != "txn1" {
		t.Errorf("FindTransaction() = %v, %v, want txn1, nil", txn, err)
	}

	_, err = p.FindTransaction("nonexistent")
	if err == nil {
		t.Error("FindTransaction() should error on non-existent transaction")
	}
}

func TestGetFirstTxn(t *testing.T) {
	p := &Plan{Txn: []transaction.Transaction{{Name: "first"}}}
	txn, err := p.GetFirstTxn()
	if err != nil || txn.Name != "first" {
		t.Errorf("GetFirstTxn() = %v, %v, want first, nil", txn, err)
	}

	empty := &Plan{}
	_, err = empty.GetFirstTxn()
	if err == nil {
		t.Error("GetFirstTxn() should error on empty plan")
	}
}
