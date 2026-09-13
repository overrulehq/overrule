# OVERRULE Core Engine
### Autonomous Patient Defense & Algorithmic Claim Reversal System

[![Go Version](https://img.shields.io/badge/go-1.24%2B-blue.svg)](https://golang.org)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPLv3-red.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Federal Statute](https://img.shields.io/badge/Statute-ERISA%20%C2%A7%20503-darkgreen.svg)](https://www.law.cornell.edu/uscode/text/29/1133)
[![Air-Gapped](https://img.shields.io/badge/Privacy-100%25%20Local--First%20(Zero--PHI)-orange.svg)](#privacy-guarantee)

OVERRULE is a sovereign, local-first legal and clinical engine built to deconstruct algorithmic health insurance claim rejections (Aetna, UnitedHealthcare, Cigna, BlueCross) and generate court-ready administrative appeal dossiers under **ERISA § 503 (29 U.S.C. § 1133)** and the **Affordable Care Act § 2719**.

---

## The Systemic Problem
Health insurers automate over 200 million claim rejections annually using black-box batch models (**Optum nH Predict**, **Cigna PXDX**, **Aetna CPB Triage**) in 1.2 seconds without meaningful physician review. Fewer than 0.2% of patients appeal due to administrative exhaustion.

OVERRULE flips the asymmetry: It converts every patient denial into a 12-page court-ready administrative record that demands internal algorithmic source code, prompt weights, and medical reviewer licensing records under **29 CFR § 2560.503-1**.

---

## Quickstart

`ash
# Clone the repository
git clone https://github.com/overrulehq/overrule.git
cd overrule

# Build binary
go build -o overrule.exe ./cmd/overrule

# List sample demonstration cases
./overrule.exe list

# Audit a specific carrier denial
./overrule.exe audit CASE-2026-AETNA-01

# Compile a certified 12-page court-ready ERISA appeal brief
./overrule.exe appeal CASE-2026-AETNA-01 -o elena_rostova_appeal.md

# Boot local-first console
./overrule.exe serve --port 5050
`

---

## Architecture & Privacy Guarantee
OVERRULE is designed under a strict **Zero-PHI Cloud Egress Guarantee**:
- 100% of text parsing, CARC code identification, and legal brief assembly execute locally on your machine.
- No patient names, diagnosis codes, or claim amounts are ever transmitted to any remote cloud API.

---

## License
Licensed under the [GNU Affero General Public License v3.0](LICENSE).
