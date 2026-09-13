package rules

import (
	"strings"

	"github.com/overrulehq/overrule/pkg/denial"
)

// PayerProfile defines known institutional denial patterns and algorithmic weaknesses for a health insurer.
type PayerProfile struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	KnownAlgorithmicTools []string `json:"known_algorithmic_tools"`
	SystemicVulnerabilities []string `json:"systemic_vulnerabilities"`
	KeyStatutoryViolations  []string `json:"key_statutory_violations"`
	AppealsAddress        string   `json:"appeals_address"`
	AppealsFax            string   `json:"appeals_fax"`
	AppealsEDIEndpoint    string   `json:"appeals_edi_endpoint"`
}

// VulnerabilityFinding represents an identified algorithmic or statutory flaw in the insurer's denial.
type VulnerabilityFinding struct {
	PayerName           string   `json:"payer_name"`
	RuleName            string   `json:"rule_name"`
	Severity            string   `json:"severity"` // "CRITICAL", "HIGH", "MEDIUM"
	Description         string   `json:"description"`
	StatutoryBasis      string   `json:"statutory_basis"`
	ReversalPrecedent   string   `json:"reversal_precedent"`
	RecommendedAction   string   `json:"recommended_action"`
}

// Catalog holds all registered payer intelligence profiles.
var Catalog = map[string]PayerProfile{
	"unitedhealthcare": {
		ID:   "uhc",
		Name: "UnitedHealthcare (Optum)",
		KnownAlgorithmicTools: []string{
			"Optum nH Predict (NaviHealth algorithmic discharge engine)",
			"Symphony AI prior-authorization triage",
			"Automated Maximum Medical Improvement (MMI) cutoff logic",
		},
		SystemicVulnerabilities: []string{
			"Algorithmic prediction model replaces individualized physician assessment of rehabilitation gains.",
			"Systematic premature termination of subacute skilled nursing facility (SNF) coverage.",
			"Denials issued under generic 'not medically necessary' boilerplate without disclosing internal benchmark parameters.",
		},
		KeyStatutoryViolations: []string{
			"ERISA § 503 (29 U.S.C. § 1133) - Failure to provide full and fair review tailored to individual patient progress.",
			"42 C.F.R. § 422.101(c) (CMS Medicare Advantage Rule) - Prohibition against using rigid internal clinical algorithms to override Medicare coverage guidelines.",
			"29 C.F.R. § 2560.503-1(g)(1)(v)(A) - Mandatory duty to provide internal clinical rule/algorithm relied upon.",
		},
		AppealsAddress:     "UnitedHealthcare Appeals & Grievances, P.O. Box 30432, Salt Lake City, UT 84130",
		AppealsFax:         "1-801-994-1083",
		AppealsEDIEndpoint: "EDI-278-UHC-DIRECT",
	},
	"aetna": {
		ID:   "aetna",
		Name: "Aetna (CVS Health)",
		KnownAlgorithmicTools: []string{
			"Clinical Policy Bulletin (CPB) Automated Matching Engine",
			"Carelon / EviCore Radiology Precertification Filter",
		},
		SystemicVulnerabilities: []string{
			"Clinical Policy Bulletins lag 18-36 months behind published medical specialty society guidelines.",
			"Improper classification of FDA-approved targeted biologics as 'experimental/investigational'.",
			"Automated denial of advanced diagnostic imaging (MRI/PET) despite documented failed conservative therapy.",
		},
		KeyStatutoryViolations: []string{
			"Affordable Care Act 45 C.F.R. § 147.136 - Internal claims and appeals requirements for evidence-based standards.",
			"29 U.S.C. § 1104(a)(1) (ERISA Fiduciary Duty) - Applying outdated policy bulletins to save payer costs at patient detriment.",
			"American College of Radiology (ACR) Appropriateness Criteria violation.",
		},
		AppealsAddress:     "Aetna Member Complaints and Appeals, P.O. Box 14463, Lexington, KY 40512",
		AppealsFax:         "1-859-455-8650",
		AppealsEDIEndpoint: "EDI-278-AETNA-CORE",
	},
	"cigna": {
		ID:   "cigna",
		Name: "Cigna Healthcare",
		KnownAlgorithmicTools: []string{
			"PXDX (Predictive Xclaim Denial Engine)",
			"EviCore automated clinical pathway filter",
		},
		SystemicVulnerabilities: []string{
			"Bulk batch denials where medical directors sign off on hundreds of rejections without opening patient charts.",
			"Arbitrary requirement of repetitive step-therapy trials for acute progressive autoimmune conditions.",
			"Failure to conduct meaningful peer-to-peer discussions within statutory timelines.",
		},
		KeyStatutoryViolations: []string{
			"State Medical Practice Act - Practicing medicine and issuing clinical necessity determinations without reviewing patient medical records.",
			"ERISA § 503 & 29 C.F.R. § 2560.503-1(h)(3)(iii) - Requirement of consultation with health care professional of appropriate training.",
			"State Prompt Payment Regulations (Failure to adjudicate within statutory limits).",
		},
		AppealsAddress:     "Cigna National Appeals Unit, P.O. Box 188011, Chattanooga, TN 37422",
		AppealsFax:         "1-877-815-4827",
		AppealsEDIEndpoint: "EDI-278-CIGNA-DIRECT",
	},
	"bluecross": {
		ID:   "bcbs",
		Name: "Blue Cross Blue Shield (Elevance / Anthem / HCSC)",
		KnownAlgorithmicTools: []string{
			"AIM Specialty Health (Carelon) Automated Authorization",
			"Step Therapy Formulary Gatekeeper",
		},
		SystemicVulnerabilities: []string{
			"Enforcing 'Fail First' protocols on patients with documented contraindications to tier-1 generic drugs.",
			"Disregarding physician documented functional decline while awaiting authorization.",
		},
		KeyStatutoryViolations: []string{
			"State Step Therapy Reform Statutes (Mandatory expedited exception upon documented clinical contraindication).",
			"Mental Health Parity and Addiction Equity Act (MHPAEA) (Non-Quantitative Treatment Limitations / NQTL violation).",
			"ERISA § 502(a)(1)(B) - Wrongful denial of vested contract benefits.",
		},
		AppealsAddress:     "Blue Cross Appeals and Grievances, P.O. Box 54159, Los Angeles, CA 90054",
		AppealsFax:         "1-888-287-7508",
		AppealsEDIEndpoint: "EDI-278-BCBS-CARE",
	},
}

// IdentifyVulnerabilities checks a denial case against known payer algorithmic patterns.
func IdentifyVulnerabilities(c *denial.DenialCase) []VulnerabilityFinding {
	findings := []VulnerabilityFinding{}
	payerKey := strings.ToLower(strings.ReplaceAll(c.InsurerName, " ", ""))

	var profile PayerProfile
	found := false
	for k, p := range Catalog {
		if strings.Contains(payerKey, k) || strings.Contains(k, payerKey) {
			profile = p
			found = true
			break
		}
	}

	if !found {
		// Generic fallback profile
		profile = PayerProfile{
			Name: c.InsurerName,
			KeyStatutoryViolations: []string{
				"ERISA § 503 (29 U.S.C. § 1133) - Failure to provide full and fair review.",
				"29 C.F.R. § 2560.503-1 - Failure to provide specific clinical rationale and guidelines.",
			},
		}
	}

	// 1. CARC Code Analysis
	if c.DenialReason.CARCCode == "CO-50" || strings.Contains(strings.ToLower(c.DenialReason.Category), "necessity") {
		findings = append(findings, VulnerabilityFinding{
			PayerName:      profile.Name,
			RuleName:       "UNSUBSTANTIATED_MEDICAL_NECESSITY_OVERRIDE",
			Severity:       "CRITICAL",
			Description:    "Payer denied coverage for clinical necessity without reconciling attending physician documentation of active impairment or disease progression.",
			StatutoryBasis: "29 C.F.R. § 2560.503-1(h)(3)(ii) requires the reviewer to not give deference to prior adverse determination and consult an independent qualified specialist.",
			ReversalPrecedent: "Federal courts routinely vacate medical necessity denials where payer's medical reviewer failed to examine the physical patient or articulate why treating physician clinical findings were discounted.",
			RecommendedAction: "Demand immediate disclosure of medical reviewer identity, specialty credentials, and internal clinical criteria under ERISA § 503.",
		})
	}

	// 2. Pre-authorization / Algorithmic Cutoff Analysis
	if c.DenialReason.CARCCode == "CO-197" || strings.Contains(strings.ToLower(c.DenialReason.Category), "authorization") {
		findings = append(findings, VulnerabilityFinding{
			PayerName:      profile.Name,
			RuleName:       "PROCEDURAL_PREAUTH_DEFENSE",
			Severity:       "HIGH",
			Description:    "Payer rejected claim on administrative grounds of lack of prior authorization, despite urgent clinical need or documented provider attempts.",
			StatutoryBasis: "ACA 45 C.F.R. § 147.136(b) - Insurers cannot deny emergency or urgent medically necessary care on strict procedural pre-service grounds.",
			ReversalPrecedent: "Administrative pre-auth denials overturned when retroactive authorization is requested with retrospective medical necessity proof.",
			RecommendedAction: "File retroactive medical necessity authorization with clinical urgency certification from attending physician.",
		})
	}

	// 3. Payer-Specific Algorithmic Exploitation
	if strings.Contains(payerKey, "united") || strings.Contains(payerKey, "optum") {
		findings = append(findings, VulnerabilityFinding{
			PayerName:      profile.Name,
			RuleName:       "OPTUM_NH_PREDICT_ALGORITHMIC_BIAS",
			Severity:       "CRITICAL",
			Description:    "UnitedHealthcare denial matches pattern of automated nH Predict cutoff without individualized clinical assessment.",
			StatutoryBasis: "CMS-4201-F & 42 C.F.R. § 422.101(c): Medicare Advantage payers cannot use internal coverage criteria that restrict access beyond traditional Medicare guidelines.",
			ReversalPrecedent: "Estate of Lokken v. UnitedHealth Group (D. Minn. 2023) demonstrated 90%+ error rate in nH Predict automated discharge targets.",
			RecommendedAction: "Incorporate Lokken legal discovery demands requiring disclosure of algorithm inputs, training weights, and actual physician review minutes.",
		})
	} else if strings.Contains(payerKey, "cigna") {
		findings = append(findings, VulnerabilityFinding{
			PayerName:      profile.Name,
			RuleName:       "PXDX_BATCH_DENIAL_CHALLENGE",
			Severity:       "CRITICAL",
			Description:    "Denial exhibits hallmarks of automated PXDX algorithm processing with sub-second medical director sign-off.",
			StatutoryBasis: "29 U.S.C. § 1104(a)(1) (ERISA Fiduciary Duty) & State Medical Board regulations on unauthorized corporate practice of medicine.",
			ReversalPrecedent: "ProPublica Cigna PXDX investigation confirmed over 300,000 claims denied without chart inspection.",
			RecommendedAction: "Subpoena medical director log records establishing whether the reviewing physician reviewed the actual imaging or diagnostic notes.",
		})
	}

	return findings
}
