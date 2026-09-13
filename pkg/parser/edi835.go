package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/overrulehq/overrule/pkg/denial"
)

// EDI835Claim represents a single claim extracted from an ANSI X12 835 Electronic Remittance Advice feed.
type EDI835Claim struct {
	ClaimID               string    `json:"claim_id"`
	StatusCode            string    `json:"status_code"`
	TotalBilled           float64   `json:"total_billed"`
	PaidAmount            float64   `json:"paid_amount"`
	PatientResponsibility float64   `json:"patient_responsibility"`
	PayerName             string    `json:"payer_name"`
	PatientLastName       string    `json:"patient_last_name"`
	PatientFirstName      string    `json:"patient_first_name"`
	MemberID              string    `json:"member_id"`
	CARCCodes             []string  `json:"carc_codes"`
	AdjustmentAmounts     []float64 `json:"adjustment_amounts"`
	ProcedureCodes        []string  `json:"procedure_codes"`
	IsDenied              bool      `json:"is_denied"`
}

// Parse835EDI processes raw ANSI X12 835 EDI content into structured claims.
func Parse835EDI(content string) ([]EDI835Claim, error) {
	claims := make([]EDI835Claim, 0)
	var currentClaim *EDI835Claim
	currentPayer := "Commercial Payer"

	delimiter := "~"
	if !strings.Contains(content, "~") && strings.Contains(content, "\n") {
		delimiter = "\n"
	}

	segments := strings.Split(content, delimiter)

	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}

		elements := strings.Split(seg, "*")
		if len(elements) == 0 {
			continue
		}

		tag := strings.TrimSpace(elements[0])

		switch tag {
		case "N1":
			if len(elements) >= 3 && elements[1] == "PR" {
				currentPayer = strings.TrimSpace(elements[2])
			}

		case "CLP":
			if currentClaim != nil && (currentClaim.IsDenied || len(currentClaim.CARCCodes) > 0) {
				claims = append(claims, *currentClaim)
			}

			claimID := ""
			status := ""
			billed := 0.0
			paid := 0.0
			ptResp := 0.0

			if len(elements) >= 2 { claimID = elements[1] }
			if len(elements) >= 3 { status = elements[2] }
			if len(elements) >= 4 { billed, _ = strconv.ParseFloat(elements[3], 64) }
			if len(elements) >= 5 { paid, _ = strconv.ParseFloat(elements[4], 64) }
			if len(elements) >= 6 { ptResp, _ = strconv.ParseFloat(elements[5], 64) }

			isDenied := status == "2" || status == "4" || paid == 0.0

			currentClaim = &EDI835Claim{
				ClaimID:               claimID,
				StatusCode:            status,
				TotalBilled:           billed,
				PaidAmount:            paid,
				PatientResponsibility: ptResp,
				PayerName:             currentPayer,
				CARCCodes:             make([]string, 0),
				AdjustmentAmounts:     make([]float64, 0),
				ProcedureCodes:        make([]string, 0),
				IsDenied:              isDenied,
			}

		case "CAS":
			if currentClaim != nil && len(elements) >= 4 {
				group := elements[1]
				code := elements[2]
				fullCode := group + "-" + code
				amt, _ := strconv.ParseFloat(elements[3], 64)

				currentClaim.CARCCodes = append(currentClaim.CARCCodes, fullCode)
				currentClaim.AdjustmentAmounts = append(currentClaim.AdjustmentAmounts, amt)
				if amt > 0 && currentClaim.PaidAmount == 0 {
					currentClaim.IsDenied = true
				}
			}

		case "NM1":
			if currentClaim != nil && len(elements) >= 5 && elements[1] == "QC" {
				currentClaim.PatientLastName = elements[3]
				currentClaim.PatientFirstName = elements[4]
				if len(elements) >= 10 && elements[8] == "MI" {
					currentClaim.MemberID = elements[9]
				}
			}

		case "SVC":
			if currentClaim != nil && len(elements) >= 2 {
				cptPart := elements[1]
				parts := strings.Split(cptPart, ":")
				if len(parts) >= 2 {
					currentClaim.ProcedureCodes = append(currentClaim.ProcedureCodes, parts[1])
				} else {
					currentClaim.ProcedureCodes = append(currentClaim.ProcedureCodes, cptPart)
				}
			}
		}
	}

	if currentClaim != nil && (currentClaim.IsDenied || len(currentClaim.CARCCodes) > 0) {
		claims = append(claims, *currentClaim)
	}

	return claims, nil
}

// ToDenialCase converts an 835 claim into a formal DenialCase.
func (c *EDI835Claim) ToDenialCase() *denial.DenialCase {
	fullName := strings.TrimSpace(c.PatientFirstName + " " + c.PatientLastName)
	if fullName == "" {
		fullName = "Claimant"
	}

	deniedAmt := c.TotalBilled - c.PaidAmount
	if deniedAmt <= 0 {
		deniedAmt = c.TotalBilled
	}

	procs := make([]denial.Procedure, 0)
	for _, pr := range c.ProcedureCodes {
		procs = append(procs, denial.Procedure{
			Code:        pr,
			Description: "Billed Service Line",
			BilledUnits: 1,
			Charge:      deniedAmt,
		})
	}
	if len(procs) == 0 {
		procs = append(procs, denial.Procedure{Code: "99214", Description: "Specialist medical intervention", BilledUnits: 1, Charge: deniedAmt})
	}

	carc := "CO-50"
	if len(c.CARCCodes) > 0 {
		carc = c.CARCCodes[0]
	}

	return &denial.DenialCase{
		ID:                    c.ClaimID,
		CreatedAt:             time.Now(),
		PatientName:           fullName,
		InsurerName:           c.PayerName,
		PolicyID:              c.MemberID,
		ClaimNumber:           c.ClaimID,
		PlanType:              denial.PlanTypeERISACommercial,
		ServiceDate:           time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
		DenialNoticeDate:      time.Now().Format("2006-01-02"),
		AppealDeadline:        time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
		TotalBilled:           c.TotalBilled,
		DeniedAmount:          deniedAmt,
		PatientResponsibility: c.PatientResponsibility,
		Diagnoses: []denial.Diagnosis{
			{Code: "R69", Description: "Illness, unspecified", Primary: true},
		},
		Procedures: procs,
		DenialReason: denial.DenialReason{
			CARCCode:      carc,
			RARCCode:      "N130",
			Category:      "Medical Necessity / Plan Limitation",
			RawLetterText: "Claim denied via 835 Remittance Advice adjudication code " + carc + ".",
		},
		Physician: denial.PhysicianInfo{
			Name:       "Attending Physician",
			Specialty:  "Clinical Specialist",
			ClinicName: "Treating Facility",
		},
		ClinicalSummary: "835 Electronic Remittance Advice automated intake.",
	}
}
