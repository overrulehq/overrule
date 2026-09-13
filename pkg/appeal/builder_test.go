package appeal

import (
	"strings"
	"testing"

	"github.com/overrulehq/overrule/pkg/denial"
)

func TestGeneratePacketDeterminism(t *testing.T) {
	c := &denial.SampleCases[0] // Elena Rostova - Aetna Brain MRI

	p1 := GeneratePacket(c)
	p2 := GeneratePacket(c)

	if p1.Case.ID != p2.Case.ID {
		t.Errorf("expected matching case IDs, got %s vs %s", p1.Case.ID, p2.Case.ID)
	}

	md1 := p1.ToMarkdown()
	md2 := p2.ToMarkdown()
	if md1 != md2 {
		t.Errorf("expected bit-stable determinism across packet generations")
	}

	if len(md1) == 0 {
		t.Fatalf("expected non-empty markdown output")
	}

	// Verify key statutory elements are present
	requiredSnippets := []string{
		"FORMAL STATUTORY & CLINICAL APPEAL",
		"ERISA Section 503",
		"29 U.S.C.",
		"29 C.F.R.",
		"TREATING PHYSICIAN CLINICAL NECESSITY DECLARATION",
		"PEER-REVIEWED SCIENTIFIC LITERATURE",
		"DISCOVERY DEMAND FOR CLAIM FILE",
	}

	for _, req := range requiredSnippets {
		if !strings.Contains(md1, req) {
			t.Errorf("markdown missing required administrative record section: %q", req)
		}
	}

	// Line count check (administrative dossier must be comprehensive)
	lines := strings.Split(md1, "\n")
	if len(lines) < 80 {
		t.Errorf("expected at least 80 lines in compiled appeal dossier, got %d", len(lines))
	}

	// Verify absence of ungrounded hyperbolic terms
	unwantedSlogans := []string{
		"federal prison",
		"criminal arrest",
		"guaranteed victory",
	}
	for _, unw := range unwantedSlogans {
		if strings.Contains(strings.ToLower(md1), unw) {
			t.Errorf("markdown contains improper hyperbolic language: %q", unw)
		}
	}
}
