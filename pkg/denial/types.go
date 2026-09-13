package denial

import (
	"time"
)

// PlanType represents the legal classification of the health plan.
type PlanType string

const (
	PlanTypeERISACommercial PlanType = "Employer Self-Funded (ERISA)"
	PlanTypeACAMarketplace  PlanType = "ACA Marketplace Individual"
	PlanTypeMedicareAdv     PlanType = "Medicare Advantage (Part C)"
	PlanTypeMedicaidManaged PlanType = "Medicaid Managed Care"
	PlanTypeTraditional     PlanType = "Traditional Commercial Fully-Insured"
)

// Diagnosis represents an ICD-10 diagnosis code and clinical description.
type Diagnosis struct {
	Code        string `json:"code"`        // e.g., "I63.9", "G35", "K50.90"
	Description string `json:"description"` // e.g., "Cerebral infarction, unspecified"
	Primary     bool   `json:"primary"`
}

// Procedure represents a CPT or HCPCS billed procedure code.
type Procedure struct {
	Code        string  `json:"code"`        // e.g., "70553", "97110", "J9271"
	Description string  `json:"description"` // e.g., "MRI brain with and without contrast"
	BilledUnits int     `json:"billed_units"`
	Charge      float64 `json:"charge"`
}

// DenialReason represents the claim adjustment codes and plain-text rationales.
type DenialReason struct {
	CARCCode    string `json:"carc_code"`    // Claim Adjustment Reason Code (e.g., "CO-50", "CO-197", "CO-16")
	RARCCode    string `json:"rarc_code"`    // Remittance Advice Remark Code (e.g., "N130", "M62")
	Category    string `json:"category"`     // "Medical Necessity", "Prior Authorization", "Experimental/Investigational", "Step Therapy"
	RawLetterText string `json:"raw_text"`   // Exact quote from insurer denial notice
}

// PhysicianInfo represents the patient's licensed treating physician.
type PhysicianInfo struct {
	Name       string `json:"name"`
	NPI        string `json:"npi"`
	Specialty  string `json:"specialty"`
	ClinicName string `json:"clinic_name"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
}

// DenialCase is the complete record of a rejected medical claim.
type DenialCase struct {
	ID                     string         `json:"id"`
	CreatedAt              time.Time      `json:"created_at"`
	PatientName            string         `json:"patient_name"`
	PatientDOB             string         `json:"patient_dob"`
	PatientState           string         `json:"patient_state"`
	InsurerName            string         `json:"insurer_name"`            // e.g. "UnitedHealthcare", "Aetna", "Cigna"
	PolicyID               string         `json:"policy_id"`
	GroupNumber            string         `json:"group_number"`
	ClaimNumber            string         `json:"claim_number"`
	PlanType               PlanType       `json:"plan_type"`
	ServiceDate            string         `json:"service_date"`
	DenialNoticeDate       string         `json:"denial_notice_date"`
	AppealDeadline         string         `json:"appeal_deadline"`
	TotalBilled            float64        `json:"total_billed"`
	DeniedAmount           float64        `json:"denied_amount"`
	PatientResponsibility  float64        `json:"patient_responsibility"`
	Diagnoses              []Diagnosis    `json:"diagnoses"`
	Procedures             []Procedure    `json:"procedures"`
	DenialReason           DenialReason   `json:"denial_reason"`
	Physician              PhysicianInfo  `json:"physician"`
	ClinicalSummary        string         `json:"clinical_summary"`        // Brief synopsis of patient condition & treatment
	PriorConservativeTried []string       `json:"prior_treatments_tried"`  // Conservative treatments already attempted
}
