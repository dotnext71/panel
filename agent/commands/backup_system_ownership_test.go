package commands

import "testing"

func TestRemapID(t *testing.T) {
	srcIDToName := map[int]string{996: "dotnext", 990: "dotnext-mail", 0: "root"}
	localNameToID := map[string]int{"dotnext": 997, "dotnext-mail": 990, "root": 0}
	tests := []struct {
		name    string
		src     int
		wantID  int
		wantHit bool
	}{
		{"dotnext shifted 996→997", 996, 997, true},
		{"dotnext-mail same id 990", 990, 0, false}, // local == src → no change
		{"root stable 0", 0, 0, false},              // local == src → no change
		{"unknown source id", 1234, 0, false},       // not in os_users.json → leave
		{"src name absent locally", 995, 0, false},  // 995 not mapped → leave
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotHit := remapID(tt.src, srcIDToName, localNameToID)
			if gotHit != tt.wantHit || (gotHit && gotID != tt.wantID) {
				t.Errorf("remapID(%d) = (%d,%v), want (%d,%v)", tt.src, gotID, gotHit, tt.wantID, tt.wantHit)
			}
		})
	}
}
