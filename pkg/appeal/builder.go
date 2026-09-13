package appeal

import (
	"fmt"
	"strings"
	"time"

	"github.com/overrulehq/overrule/pkg/clinical"
	"github.com/overrulehq/overrule/pkg/denial"
	"github.com/overrulehq/overrule/pkg/legal"
	"github.com/overrulehq/overrule/pkg/rules"
)

// AppealPacket contains the complete clinical and statutory appeal dossier.
type AppealPacket struct {
	Case                *denial.DenialCase           `json:"case"`
	GeneratedAt         time.Time                    `json:"generated_at"`
	Vulnerabilities     []rules.VulnerabilityFinding `json:"vulnerabilities"`
	ClinicalArguments   []clinical.ClinicalArgument  `json:"clinical_arguments"`
	PhysicianRebuttal   string                       `json:"physician_rebuttal"`
	LegalBrief          legal.LegalBrief             `json:"legal_brief"`
	SubmissionChecklist []string                     `json:"submission_checklist"`
}

// GeneratePacket processes a DenialCase and compiles the comprehensive appeal.
func GeneratePacket(c *denial.DenialCase) *AppealPacket {
	vulnerabilities := rules.IdentifyVulnerabilities(c)
	clinicalArgs := clinical.BuildClinicalArguments(c)
	physicianStatement := clinical.GeneratePhysicianRebuttal(c, clinicalArgs)
	legalBrief := legal.BuildLegalBrief(c, vulnerabilities)

	checklist := []string{
		"1. Print or export this certified Overrule Appeal Dossier in full.",
		fmt.Sprintf("2. Have attending physician (%s, MD) sign Section III or attach signed letter.", c.Physician.Name),
		"3. Attach Exhibit A: Complete copy of insurer's original denial letter / EOB.",
		"4. Attach Exhibit B: Relevant medical chart notes, imaging reports, and lab findings.",
		"5. Attach Exhibit C: Attached peer-reviewed literature prints (PubMed citations).",
		fmt.Sprintf("6. Transmit via Certified Mail with Return Receipt Requested (USPS Form 3800) to %s.", c.InsurerName),
		"7. Transmit digital copy via insurer member portal appeals upload tool.",
		"8. File concurrent notice with State Insurance Commissioner Consumer Complaint Portal.",
	}

	return &AppealPacket{
		Case:                c,
		GeneratedAt:         time.Now(),
		Vulnerabilities:     vulnerabilities,
		ClinicalArguments:   clinicalArgs,
		PhysicianRebuttal:   physicianStatement,
		LegalBrief:          legalBrief,
		SubmissionChecklist: checklist,
	}
}

// ToMarkdown produces a formal, court-ready Markdown document of the appeal.
func (p *AppealPacket) ToMarkdown() string {
	c := p.Case
	sb := strings.Builder{}

	sb.WriteString("# FORMAL STATUTORY & CLINICAL APPEAL OF ADVERSE BENEFIT DETERMINATION\n")
	sb.WriteString("### Prepared via Overrule Autonomous Patient Defense System\n\n")
	sb.WriteString("---\n\n")

	sb.WriteString("## I. CASE IDENTIFICATION & TIMELINE\n\n")
	sb.WriteString(fmt.Sprintf("**Date of Appeal:** %s\n\n", p.GeneratedAt.Format("January 02, 2006")))
	sb.WriteString("| Field | Details |\n| :--- | :--- |\n")
	sb.WriteString(fmt.Sprintf("| **Patient / Claimant** | %s |\n", c.PatientName))
	sb.WriteString(fmt.Sprintf("| **Date of Birth** | %s |\n", c.PatientDOB))
	sb.WriteString(fmt.Sprintf("| **Patient State of Residence** | %s |\n", c.PatientState))
	sb.WriteString(fmt.Sprintf("| **Health Insurance Carrier** | %s |\n", c.InsurerName))
	sb.WriteString(fmt.Sprintf("| **Policy / Member ID** | %s |\n", c.PolicyID))
	sb.WriteString(fmt.Sprintf("| **Group / Account Number** | %s |\n", c.GroupNumber))
	sb.WriteString(fmt.Sprintf("| **Adverse Claim Number** | %s |\n", c.ClaimNumber))
	sb.WriteString(fmt.Sprintf("| **Plan Legal Classification** | %s |\n", c.PlanType))
	sb.WriteString(fmt.Sprintf("| **Date of Service** | %s |\n", c.ServiceDate))
	sb.WriteString(fmt.Sprintf("| **Adverse Notice Date** | %s |\n", c.DenialNoticeDate))
	sb.WriteString(fmt.Sprintf("| **Total Billed Amount** | $%.2f |\n", c.TotalBilled))
	sb.WriteString(fmt.Sprintf("| **Improperly Denied Balance** | $%.2f |\n", c.DeniedAmount))
	sb.WriteString(fmt.Sprintf("| **Alleged Patient Responsibility** | $%.2f |\n\n", c.PatientResponsibility))

	sb.WriteString("### Diagnoses (ICD-10):\n")
	for _, d := range c.Diagnoses {
		sb.WriteString(fmt.Sprintf("- **%s**: %s\n", d.Code, d.Description))
	}
	sb.WriteString("\n### Disallowed Procedures / Services (CPT / HCPCS):\n")
	for _, pr := range c.Procedures {
		sb.WriteString(fmt.Sprintf("- **%s**: %s (Billed: $%.2f)\n", pr.Code, pr.Description, pr.Charge))
	}
	sb.WriteString(fmt.Sprintf("\n**Purported Denial Rationale (Code: %s):**\n> *\"%s\"*\n\n", c.DenialReason.CARCCode, c.DenialReason.RawLetterText))

	sb.WriteString("---\n\n")
	sb.WriteString("## II. EXECUTIVE SUMMARY & SYSTEMIC AUDIT FINDINGS\n\n")
	sb.WriteString(fmt.Sprintf("Claimant hereby formally appeals the adverse determination rendered by %s on %s. A thorough multi-disciplinary audit of the medical record and insurance policy reveals that the denial is clinically unsound, contradicts established medical specialty standards, and violates federal statutory claims regulations under ERISA Section 503 and the Affordable Care Act.\n\n", c.InsurerName, c.DenialNoticeDate))

	if len(p.Vulnerabilities) > 0 {
		sb.WriteString("### Identified Institutional & Algorithmic Flaws:\n")
		for _, v := range p.Vulnerabilities {
			sb.WriteString(fmt.Sprintf("#### [%s] %s\n", v.Severity, v.RuleName))
			sb.WriteString(fmt.Sprintf("- **Finding:** %s\n", v.Description))
			sb.WriteString(fmt.Sprintf("- **Statutory Basis:** %s\n", v.StatutoryBasis))
			sb.WriteString(fmt.Sprintf("- **Judicial / Regulatory Precedent:** %s\n\n", v.ReversalPrecedent))
		}
	}

	sb.WriteString("---\n\n")
	sb.WriteString("## III. TREATING PHYSICIAN CLINICAL NECESSITY DECLARATION\n\n")
	sb.WriteString("```text\n")
	sb.WriteString(p.PhysicianRebuttal)
	sb.WriteString("\n```\n\n")

	sb.WriteString("---\n\n")
	sb.WriteString("## IV. PEER-REVIEWED SCIENTIFIC LITERATURE & GUIDELINE EXHIBITS\n\n")
	for _, ca := range p.ClinicalArguments {
		sb.WriteString(fmt.Sprintf("### Clinical Authority for %s (%s)\n", ca.ProcedureCode, ca.DiagnosisCode))
		sb.WriteString(fmt.Sprintf("**Specialty Standard of Care:** %s\n\n", ca.StandardOfCare))
		for _, lit := range ca.Literature {
			sb.WriteString(fmt.Sprintf("* **%s**\n", lit.Citation))
			sb.WriteString(fmt.Sprintf("  * *Journal:* %s | *PMID:* [%s](https://pubmed.ncbi.nlm.nih.gov/%s/) | *Grade:* %s\n", lit.Journal, lit.PMID, lit.PMID, lit.Grade))
			sb.WriteString(fmt.Sprintf("  * *Key Clinical Finding:* %s\n\n", lit.KeyFinding))
		}
	}

	sb.WriteString("---\n\n")
	sb.WriteString("## V. STATUTORY LEGAL FOUNDATIONS & REGULATORY VIOLATIONS\n\n")
	for _, sec := range p.LegalBrief.Sections {
		sb.WriteString(fmt.Sprintf("### %s\n", sec.Title))
		sb.WriteString(fmt.Sprintf("**Legal Citation:** `%s`\n\n", sec.Citation))
		sb.WriteString(fmt.Sprintf("**Specific Violation:**\n%s\n\n", sec.Violation))
		sb.WriteString(fmt.Sprintf("**Required Legal Remedy:**\n%s\n\n", sec.Remedy))
	}

	sb.WriteString("### VI. FORMAL DISCOVERY DEMAND FOR CLAIM FILE\n\n")
	sb.WriteString("```text\n")
	sb.WriteString(p.LegalBrief.DocumentDemandNotice)
	sb.WriteString("\n```\n\n")

	sb.WriteString("### VII. BAD FAITH CLAIM SETTLEMENT WARNING\n\n")
	sb.WriteString("```text\n")
	sb.WriteString(p.LegalBrief.BadFaithWarning)
	sb.WriteString("\n```\n\n")

	sb.WriteString("### VIII. REGULATORY ESCALATION NOTICE\n\n")
	sb.WriteString("```text\n")
	sb.WriteString(p.LegalBrief.RegulatoryEscalationNotice)
	sb.WriteString("\n```\n\n")

	sb.WriteString("---\n\n")
	sb.WriteString("## IX. FILING CHECKLIST & TRANSMISSION PROTOCOL\n\n")
	for _, item := range p.SubmissionChecklist {
		sb.WriteString(fmt.Sprintf("- [ ] %s\n", item))
	}

	sb.WriteString("\n---\n")
	sb.WriteString("*Generated by Overrule Core (v1.0.0-beta) — Built to Protect Patients & Eliminate Algorithmic Healthcare Injustice.*\n")

	return sb.String()
}
