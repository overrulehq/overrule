package legal

import (
	"fmt"

	"github.com/overrulehq/overrule/pkg/denial"
	"github.com/overrulehq/overrule/pkg/rules"
)

// StatutorySection represents a specific legal authority invoked in the appeal.
type StatutorySection struct {
	Title       string `json:"title"`
	Citation    string `json:"citation"`
	Violation   string `json:"violation"`
	Remedy      string `json:"remedy"`
}

// LegalBrief holds the complete statutory foundation of the appeal.
type LegalBrief struct {
	Sections               []StatutorySection `json:"sections"`
	DocumentDemandNotice   string             `json:"document_demand_notice"`
	BadFaithWarning        string             `json:"bad_faith_warning"`
	RegulatoryEscalationNotice string         `json:"regulatory_escalation_notice"`
}

// BuildLegalBrief constructs an airtight legal brief tailored to the plan type and payer.
func BuildLegalBrief(c *denial.DenialCase, findings []rules.VulnerabilityFinding) LegalBrief {
	sections := []StatutorySection{}

	// 1. ERISA § 503 (29 U.S.C. § 1133) & 29 C.F.R. § 2560.503-1
	sections = append(sections, StatutorySection{
		Title:    "Denial of Full and Fair Review Under Federal Law",
		Citation: "ERISA Section 503 (29 U.S.C. § 1133) & 29 C.F.R. § 2560.503-1(h)",
		Violation: fmt.Sprintf(
			"%s issued an adverse benefit determination without affording the claimant a full and fair review. Federal regulations mandate that review on appeal must not afford deference to the initial adverse determination and must be conducted by an appropriate named fiduciary who is neither the individual who made the adverse determination nor the subordinate of such individual. In issuing a boilerplate denial for %s, the plan failed to articulate individualized, clinical grounds.",
			c.InsurerName, c.PatientName,
		),
		Remedy: "Immediate de novo clinical review by an independent, board-certified physician in the relevant medical specialty who has no financial conflict of interest.",
	})

	// 2. Mandatory Production of Internal Criteria (29 C.F.R. § 2560.503-1(g)(1)(v)(A))
	sections = append(sections, StatutorySection{
		Title:    "Failure to Disclose Internal Clinical Guidelines and Algorithmic Rules",
		Citation: "29 C.F.R. § 2560.503-1(g)(1)(v)(A)",
		Violation: fmt.Sprintf(
			"The plan's denial notice references vague 'medical necessity' benchmarks or proprietary coverage policies without providing the specific internal rules, guidelines, protocols, or algorithmic criteria relied upon in making the adverse determination. Under federal law, if an adverse benefit determination is based on an internal rule or protocol, the notice must provide either the specific rule or a statement that such rule will be provided free of charge upon request.",
			),
		Remedy: "Immediate disclosure of the full clinical criteria, software algorithms, or clinical policy bulletins utilized.",
	})

	// 3. Affordable Care Act Protections (45 C.F.R. § 147.136)
	sections = append(sections, StatutorySection{
		Title:    "Affordable Care Act Internal Appeals and External Review Rights",
		Citation: "Public Health Service Act § 2719 (42 U.S.C. § 300gg-19) & 45 C.F.R. § 147.136",
		Violation: fmt.Sprintf(
			"Under the ACA, non-grandfathered group health plans and health insurance issuers must provide an internal claims and appeals process that complies with the Department of Labor claims procedure, and must permit the claimant to review the claim file and to present evidence and testimony. Continued denial without consideration of the attached peer-reviewed medical evidence entitles claimant to immediate binding External Review by an Independent Review Organization (IRO).",
		),
		Remedy: "Certification of claim for expedited external review if internal appeal is not reversed within statutory deadlines.",
	})

	// 4. Incorporate Payer-Specific Vulnerabilities
	for _, f := range findings {
		sections = append(sections, StatutorySection{
			Title:     fmt.Sprintf("Payer Systemic Violation: %s", f.RuleName),
			Citation:  f.StatutoryBasis,
			Violation: f.Description,
			Remedy:    f.RecommendedAction,
		})
	}

	// Document Demand under ERISA § 104(b)(4) & 29 C.F.R. § 2560.503-1(h)(2)(iii)
	docDemand := fmt.Sprintf(`FORMAL STATUTORY DEMAND FOR PRODUCTION OF COMPLETE CLAIM FILE:
Pursuant to 29 C.F.R. § 2560.503-1(h)(2)(iii) and ERISA Section 104(b)(4) (29 U.S.C. § 1024(b)(4)), claimant hereby formally requests copies of all documents, records, and other information relevant to this claim for benefits, free of charge.

This demand includes, but is not limited to:
1. The complete administrative claim file, including all medical records, notes, and correspondence.
2. The specific internal rules, guidelines, clinical protocols, or algorithmic scoring tools relied upon.
3. The names, professional medical credentials, state medical license numbers, and board certifications of all medical directors, consultants, or reviewers who participated in the adverse determination.
4. The minutes of any peer-to-peer or internal clinical review panel meetings regarding this claim.
5. Notes and time logs establishing the exact duration spent by the medical director reviewing the actual patient chart.

NOTICE: Under 29 U.S.C. § 1132(c)(1), failure by a plan administrator to furnish requested documents within thirty (30) days may subject the plan administrator to statutory personal liability penalties of up to $110 per day.`)

	badFaith := fmt.Sprintf(`BAD FAITH CLAIM SETTLEMENT WARNING:
Please be advised that systematic denial of medically necessary services, refusal to engage in meaningful dialogue with treating physicians, or reliance on automated algorithmic screening tools that override individual patient clinical records constitutes bad faith insurance practices under State law and a breach of fiduciary duty under ERISA Section 404(a)(1) (29 U.S.C. § 1104(a)(1)). 

Claimant reserves all rights to seek full benefit recovery, prejudgment interest, statutory penalties, and reasonable attorney's fees pursuant to 29 U.S.C. § 1132(g)(1).`)

	regulatory := fmt.Sprintf(`REGULATORY ESCALATION CONCURRENT NOTICE:
Copies of this appeal, along with the complete evidentiary record and proof of certified delivery, are being lodged with:
1. The State Insurance Commissioner / Department of Managed Health Care (Consumer Assistance Unit).
2. The United States Department of Labor - Employee Benefits Security Administration (EBSA).
3. The Centers for Medicare & Medicaid Services (CMS) Oversight Division (if applicable).

If this claim is not fully reversed within the statutory 30-day window (or 72 hours if certified as clinically urgent), an immediate complaint for unfair settlement practices and request for independent external audit will be prosecuted.`)

	return LegalBrief{
		Sections:                   sections,
		DocumentDemandNotice:       docDemand,
		BadFaithWarning:            badFaith,
		RegulatoryEscalationNotice: regulatory,
	}
}
