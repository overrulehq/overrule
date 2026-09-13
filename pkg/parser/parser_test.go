package parser

import (
	"testing"
)

func TestParseRawEOBText(t *testing.T) {
	sampleMessyOCR := `
	EXPLANATION OF BENEFITS - THIS IS NOT A BILL
	UNITEDHEALTHCARE SERVICES, INC.
	Payer ID: 87726 | Date: 08/14/2026
	
	Subscriber Name: HENDERSON, ROBERT M
	Member ID: UHC-88129401-01
	Group Number: GRP-EXEC-9021
	Claim Number: CLM-2026-99412A
	
	Service Date: 07/22/2026
	Provider: ST. JUDE SUBACUTE NEUROLOGICAL REHAB
	
	Billed Services:
	CPT 97110 - Therapeutic exercises 15 min - Billed: $1,200.00
	CPT 97530 - Therapeutic activities direct patient contact - Billed: $2,400.00
	
	Adjudication Detail:
	Total Billed: $3,600.00
	Paid: $0.00
	Patient Responsibility: $3,600.00
	
	Claim Adjustment Reason Codes:
	CO-50: These are non-covered services because this is not deemed a medical necessity.
	Remark Code: N130
	
	Reason for Adverse Determination:
	Services rendered exceed plan coverage parameters for post-acute subacute rehabilitation.
	Maximum Medical Improvement reached under internal clinical utilization model.
	`

	pd := ParseRawText(sampleMessyOCR)

	if pd.PayerName != "UnitedHealthcare" {
		t.Errorf("expected PayerName UnitedHealthcare, got %s", pd.PayerName)
	}
	if pd.ClaimNumber != "CLM-2026-99412A" {
		t.Errorf("expected ClaimNumber CLM-2026-99412A, got %s", pd.ClaimNumber)
	}
	if pd.PolicyID != "UHC-88129401-01" {
		t.Errorf("expected PolicyID UHC-88129401-01, got %s", pd.PolicyID)
	}
	if len(pd.CARCCodes) == 0 || pd.CARCCodes[0] != "CO-50" {
		t.Errorf("expected CARC code CO-50, got %v", pd.CARCCodes)
	}
	if pd.TotalBilled != 3600.00 {
		t.Errorf("expected TotalBilled 3600.00, got %.2f", pd.TotalBilled)
	}
}

func TestParse835EDI(t *testing.T) {
	raw835 := "ISA*00*          *00*          *ZZ*SUBMITTER      *ZZ*RECEIVER       *260815*1030*U*00501*000000001*0*T*:~GS*HP*SUBMITTER*RECEIVER*20260815*1030*1*X*005010X221A1~BPR*I*0.00*C*ACH*CTX*01*999999999*DA*12345678*1999999999**01*999999999*DA*87654321*20260815~N1*PR*AETNA LIFE INSURANCE~CLP*CLM-AETNA-8812*2*4250.00*0.00*4250.00*12*AETNA9912001**1~CAS*CO*50*4250.00*1~NM1*QC*1*ROSTOVA*ELENA****MI*AET-90412891~SVC*HC:70553*4250.00*0.00~SE*8*0001~GE*1*1~IEA*1*000000001~"

	claims, err := Parse835EDI(raw835)
	if err != nil {
		t.Fatalf("unexpected error parsing 835: %v", err)
	}

	if len(claims) != 1 {
		t.Fatalf("expected 1 claim, got %d", len(claims))
	}

	c := claims[0]
	if c.ClaimID != "CLM-AETNA-8812" {
		t.Errorf("expected ClaimID CLM-AETNA-8812, got %s", c.ClaimID)
	}
	if !c.IsDenied {
		t.Errorf("expected claim to be marked denied")
	}
	if c.PayerName != "AETNA LIFE INSURANCE" {
		t.Errorf("expected PayerName AETNA LIFE INSURANCE, got %s", c.PayerName)
	}
	if len(c.CARCCodes) == 0 || c.CARCCodes[0] != "CO-50" {
		t.Errorf("expected CARC CO-50, got %v", c.CARCCodes)
	}
	if c.PatientLastName != "ROSTOVA" || c.PatientFirstName != "ELENA" {
		t.Errorf("expected patient Elena Rostova, got %s %s", c.PatientFirstName, c.PatientLastName)
	}
}
