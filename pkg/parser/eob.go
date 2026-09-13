package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/overrulehq/overrule/pkg/denial"
)

// RawDenialInput represents unstructured text extracted from an EOB, denial notice, or OCR scan.
type RawDenialInput struct {
	Text     string    `json:"text"`
	Source   string    `json:"source"` // "ocr_text", "pasted_text", "pdf_extract"
	Filename string    `json:"filename,omitempty"`
	ScanDate time.Time `json:"scan_date"`
}

// ParsedDenial contains the normalized clinical, legal, and financial entities extracted from raw text.
type ParsedDenial struct {
	PayerName             string   `json:"payer_name"`
	ClaimNumber           string   `json:"claim_number"`
	PolicyID              string   `json:"policy_id"`
	GroupNumber           string   `json:"group_number"`
	PatientName           string   `json:"patient_name"`
	PatientDOB            string   `json:"patient_dob"`
	ServiceDate           string   `json:"service_date"`
	DenialDate            string   `json:"denial_date"`
	AppealDeadline        string   `json:"appeal_deadline"`
	TotalBilled           float64  `json:"total_billed"`
	DeniedAmount          float64  `json:"denied_amount"`
	PatientResponsibility float64  `json:"patient_responsibility"`
	CARCCodes             []string `json:"carc_codes"`
	RARCCodes             []string `json:"rarc_codes"`
	DiagnosisCodes        []string `json:"diagnosis_codes"`
	ProcedureCodes        []string `json:"procedure_codes"`
	RawDenialExcerpt      string   `json:"raw_denial_excerpt"`
	ConfidenceScore       float64  `json:"confidence_score"`
}

var (
	reClaimNum = regexp.MustCompile(`(?i)(?:claim\s*(?:#|no|number|id)?[:\s]+)([A-Z0-9\-]{6,24})`)
	rePolicyID = regexp.MustCompile(`(?i)(?:member\s*(?:id|#|no)?|policy\s*(?:id|#|no)?|subscriber\s*id)[:\s]+([A-Z0-9\-]{6,20})`)
	reGroupNum = regexp.MustCompile(`(?i)(?:group\s*(?:#|no|number)?[:\s]+)([A-Z0-9\-]{4,16})`)
	reCARC     = regexp.MustCompile(`\b(CO|PR|OA|PI|CR)-?(\d{1,3})\b`)
	reRARC     = regexp.MustCompile(`\b([NMN]\d{2,4})\b`)
	reICD10    = regexp.MustCompile(`\b([A-TV-Z]\d{2}(?:\.\d{1,4})?)\b`)
	reCPT      = regexp.MustCompile(`\b(\d{4}[0-9A-Z]|[A-V]\d{4})\b`)
	reMoney    = regexp.MustCompile(`\$\s*([0-9]{1,3}(?:,[0-9]{3})*(?:\.[0-9]{2})?)`)
	reDate     = regexp.MustCompile(`\b(\d{1,2}[/-]\d{1,2}[/-]\d{2,4}|\b(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\s+\d{1,2},?\s+\d{4})\b`)
)

// ParseRawText extracts structured denial fields from raw EOB/notice text.
func ParseRawText(raw string) *ParsedDenial {
	pd := &ParsedDenial{
		CARCCodes:      make([]string, 0),
		RARCCodes:      make([]string, 0),
		DiagnosisCodes: make([]string, 0),
		ProcedureCodes: make([]string, 0),
	}

	lower := strings.ToLower(raw)
	switch {
	case strings.Contains(lower, "unitedhealthcare") || strings.Contains(lower, "united healthcare") || strings.Contains(lower, "optum"):
		pd.PayerName = "UnitedHealthcare"
	case strings.Contains(lower, "aetna") || strings.Contains(lower, "cvs health"):
		pd.PayerName = "Aetna"
	case strings.Contains(lower, "cigna") || strings.Contains(lower, "evernorth"):
		pd.PayerName = "Cigna Healthcare"
	case strings.Contains(lower, "blue cross") || strings.Contains(lower, "blue shield") || strings.Contains(lower, "anthem") || strings.Contains(lower, "bcbs"):
		pd.PayerName = "Blue Cross Blue Shield"
	case strings.Contains(lower, "humana"):
		pd.PayerName = "Humana"
	case strings.Contains(lower, "kaiser"):
		pd.PayerName = "Kaiser Permanente"
	default:
		pd.PayerName = "Commercial Carrier"
	}

	if m := reClaimNum.FindStringSubmatch(raw); len(m) > 1 {
		pd.ClaimNumber = strings.TrimSpace(m[1])
	}
	if m := rePolicyID.FindStringSubmatch(raw); len(m) > 1 {
		pd.PolicyID = strings.TrimSpace(m[1])
	}
	if m := reGroupNum.FindStringSubmatch(raw); len(m) > 1 {
		pd.GroupNumber = strings.TrimSpace(m[1])
	}

	for _, match := range reCARC.FindAllStringSubmatch(raw, -1) {
		code := match[1] + "-" + match[2]
		if !contains(pd.CARCCodes, code) {
			pd.CARCCodes = append(pd.CARCCodes, code)
		}
	}

	for _, match := range reRARC.FindAllStringSubmatch(raw, -1) {
		code := match[1]
		if !contains(pd.RARCCodes, code) {
			pd.RARCCodes = append(pd.RARCCodes, code)
		}
	}

	for _, match := range reICD10.FindAllString(raw, -1) {
		if !strings.HasPrefix(match, "CO") && !strings.HasPrefix(match, "PR") && len(match) >= 3 {
			if !contains(pd.DiagnosisCodes, match) {
				pd.DiagnosisCodes = append(pd.DiagnosisCodes, match)
			}
		}
	}

	for _, match := range reCPT.FindAllString(raw, -1) {
		if match != "2024" && match != "2025" && match != "2026" && match != "2027" {
			if !contains(pd.ProcedureCodes, match) {
				pd.ProcedureCodes = append(pd.ProcedureCodes, match)
			}
		}
	}

	moneyMatches := reMoney.FindAllStringSubmatch(raw, -1)
	amounts := make([]float64, 0)
	for _, m := range moneyMatches {
		clean := strings.ReplaceAll(m[1], ",", "")
		if val, err := strconv.ParseFloat(clean, 64); err == nil && val > 0 {
			amounts = append(amounts, val)
		}
	}
	if len(amounts) > 0 {
		maxAmt := 0.0
		for _, a := range amounts {
			if a > maxAmt {
				maxAmt = a
			}
		}
		pd.TotalBilled = maxAmt
		pd.DeniedAmount = maxAmt
	}

	dateMatches := reDate.FindAllString(raw, -1)
	if len(dateMatches) > 0 {
		pd.ServiceDate = dateMatches[0]
	}
	if len(dateMatches) > 1 {
		pd.DenialDate = dateMatches[1]
	}

	lines := strings.Split(raw, "\n")
	excerptLines := make([]string, 0)
	for _, l := range lines {
		tl := strings.TrimSpace(l)
		ll := strings.ToLower(tl)
		if strings.Contains(ll, "denied") || strings.Contains(ll, "not medically necessary") ||
			strings.Contains(ll, "not covered") || strings.Contains(ll, "prior authorization") ||
			strings.Contains(ll, "criteria") || strings.Contains(ll, "policy bulletin") {
			excerptLines = append(excerptLines, tl)
		}
	}
	if len(excerptLines) > 0 {
		pd.RawDenialExcerpt = strings.Join(excerptLines, " ")
	}

	score := 0.0
	if pd.PayerName != "Commercial Carrier" { score += 0.25 }
	if pd.ClaimNumber != "" { score += 0.20 }
	if len(pd.CARCCodes) > 0 { score += 0.25 }
	if pd.DeniedAmount > 0 { score += 0.15 }
	if len(pd.DiagnosisCodes) > 0 || len(pd.ProcedureCodes) > 0 { score += 0.15 }
	pd.ConfidenceScore = score

	return pd
}

// ToDenialCase converts the parsed raw data into the system DenialCase domain model.
func (p *ParsedDenial) ToDenialCase(patientName string, physicianName string) *denial.DenialCase {
	if patientName == "" {
		patientName = p.PatientName
		if patientName == "" {
			patientName = "Claimant"
		}
	}

	diag := make([]denial.Diagnosis, 0)
	for i, d := range p.DiagnosisCodes {
		diag = append(diag, denial.Diagnosis{
			Code:        d,
			Description: "Extracted diagnostic finding",
			Primary:     i == 0,
		})
	}
	if len(diag) == 0 {
		diag = append(diag, denial.Diagnosis{Code: "R69", Description: "Illness, unspecified", Primary: true})
	}

	proc := make([]denial.Procedure, 0)
	for _, pr := range p.ProcedureCodes {
		proc = append(proc, denial.Procedure{
			Code:        pr,
			Description: "Extracted medical service",
			BilledUnits: 1,
			Charge:      p.DeniedAmount,
		})
	}
	if len(proc) == 0 {
		proc = append(proc, denial.Procedure{Code: "99214", Description: "Specialist consultation", BilledUnits: 1, Charge: p.DeniedAmount})
	}

	carc := "CO-50"
	if len(p.CARCCodes) > 0 {
		carc = p.CARCCodes[0]
	}

	rarc := "N130"
	if len(p.RARCCodes) > 0 {
		rarc = p.RARCCodes[0]
	}

	id := p.ClaimNumber
	if id == "" {
		id = "CLAIM-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	serviceDate := p.ServiceDate
	if serviceDate == "" {
		serviceDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}

	denialDate := p.DenialDate
	if denialDate == "" {
		denialDate = time.Now().Format("2006-01-02")
	}

	return &denial.DenialCase{
		ID:                    id,
		CreatedAt:             time.Now(),
		PatientName:           patientName,
		InsurerName:           p.PayerName,
		PolicyID:              p.PolicyID,
		GroupNumber:           p.GroupNumber,
		ClaimNumber:           p.ClaimNumber,
		PlanType:              denial.PlanTypeERISACommercial,
		ServiceDate:           serviceDate,
		DenialNoticeDate:      denialDate,
		AppealDeadline:        time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
		TotalBilled:           p.TotalBilled,
		DeniedAmount:          p.DeniedAmount,
		PatientResponsibility: p.DeniedAmount,
		Diagnoses:             diag,
		Procedures:            proc,
		DenialReason: denial.DenialReason{
			CARCCode:      carc,
			RARCCode:      rarc,
			Category:      "Medical Necessity / Adverse Determination",
			RawLetterText: p.RawDenialExcerpt,
		},
		Physician: denial.PhysicianInfo{
			Name:       physicianName,
			Specialty:  "Treating Clinical Specialist",
			ClinicName: "Attending Medical Facility",
		},
		ClinicalSummary: "Contested under ERISA § 503 and governing clinical specialty guidelines.",
	}
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
