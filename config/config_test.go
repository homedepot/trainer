package config

import (
"testing"

"github.com/homedepot/trainer/structs/plan"
)

func TestNewConfig_RejectsPathTraversal(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"path traversal", "../../../etc/passwd", true},
		{"double dot in middle", "data/../../secrets.yml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
_, err := NewConfig(tt.filename, false, "", nil)
if (err != nil) != tt.wantErr {
t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
}
})
	}
}

func TestLoadVars_RejectsPathTraversal(t *testing.T) {
	cfg := &Config{
		Plans: []plan.Plan{{Name: "test", ExtVarFile: "../../etc/passwd"}},
	}

	err := cfg.LoadVars()
	if err == nil {
		t.Error("LoadVars() should reject path traversal")
	}
}

func TestFindPlan(t *testing.T) {
	cfg := &Config{
		Plans: []plan.Plan{{Name: "plan1"}, {Name: "plan2"}},
	}

	p, err := cfg.FindPlan("plan1")
	if err != nil || p.Name != "plan1" {
		t.Errorf("FindPlan() = %v, %v, want plan1, nil", p, err)
	}

	_, err = cfg.FindPlan("nonexistent")
	if err == nil {
		t.Error("FindPlan() should error on non-existent plan")
	}
}
