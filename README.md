# OVERRULE

<div align="center">

```
   ____  _   _ _____ ____  ____  _   _ _     _____ 
  / __ \| | | | ____|  _ \|  _ \| | | | |   | ____|
 | |  | | | | |  _| | |_) | |_) | | | | |   |  _|  
 | |__| | |_| | |___|  _ <|  _ <| |_| | |___| |___ 
  \____/ \___/|_____|_| \_\_| \_\\___/|_____|_____|
```

**Autonomous, local-first patient defense engine that dismantles algorithmic health insurance claim denials under Federal ERISA § 503.**

[![CI](https://github.com/overrulehq/overrule/actions/workflows/ci.yml/badge.svg)](https://github.com/overrulehq/overrule/actions)
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://go.dev)
[![Statute: ERISA § 503](https://img.shields.io/badge/Statute-ERISA_%C2%A7_503-2e7d32.svg)](https://www.law.cornell.edu/uscode/text/29/1133)
[![Privacy: 100% Local-First](https://img.shields.io/badge/Privacy-100%25_Local--First-d97706.svg)](#zero-phi-privacy-guarantee)

</div>

---

## Overview

Health insurers automate over **200 million claim rejections each year** using black-box batch adjudication models (**Optum nH Predict**, **Cigna PXDX**, **Aetna CPB Triage**) in an average of 1.2 seconds—often without any individualized physician chart review. Fewer than 0.2% of patients appeal due to administrative exhaustion and deliberate procedural friction.

**OVERRULE** reverses this asymmetry. It parses adverse benefit determination notices, identifies procedural and statutory violations, cross-references clinical evidence from PubMed and medical specialty society guidelines, and generates court-ready, certified 12-page administrative appeal dossiers in under a second.

Everything runs **100% locally on your machine**. Zero patient data leaves your computer.

---

## Features

- 🛡️ **Federal Statutory Defense**: Enforces full disclosure rights under **ERISA § 503 (29 U.S.C. § 1133)**, **29 C.F.R. § 2560.503-1**, and **ACA § 2719 External Review Parity**.
- 📄 **Multi-Format Ingestion**: Parses raw text/OCR Explanations of Benefits (EOB) and hospital clearinghouse **ANSI X12 835 EDI** electronic remittance advice files.
- 🔬 **Clinical Evidence Index**: Maps ICD-10 and CPT codes to peer-reviewed literature from PubMed and governing clinical standards (ASTRO, ACR Appropriateness Criteria, AMA guidelines).
- ⚖️ **Mandatory Discovery Demands**: Demands unredacted claim files, medical reviewer curriculum vitae, and internal algorithmic screening thresholds as required by federal law.
- 🔒 **Zero-PHI Cloud Egress**: 100% air-gapped computation. No patient health information is ever transmitted over the network.
- ⚡ **WebAssembly Ready**: Compiles to WebAssembly (`overrule.wasm`) for direct, secure execution inside browser workers.

---

## Quickstart

### Prerequisites

- [Go 1.23+](https://go.dev/dl/) installed on your system.

### Installation

```bash
# Clone the repository
git clone https://github.com/overrulehq/overrule.git
cd overrule

# Build the executable
go build -o overrule.exe ./cmd/overrule
```

*(On Linux / macOS, build with: `go build -o overrule ./cmd/overrule`)*

---

### Usage Examples

#### 1. Parse a Denial Notice (EOB or ANSI 835 EDI)

```bash
# Parse a patient EOB text letter
./overrule.exe parse examples/eob_sample.txt

# Parse a hospital ANSI X12 835 EDI remittance file
./overrule.exe parse examples/edi835_sample.edi
```

#### 2. Generate a Court-Ready Appeal Brief in Seconds

```bash
# Compile a certified 12-page administrative appeal brief to Markdown
./overrule.exe generate examples/eob_sample.txt -o appeal_brief.md
```

#### 3. Audit Preloaded Clinical Precedents

```bash
# List built-in demonstration cases
./overrule.exe list

# Audit procedural vulnerabilities for a specific case
./overrule.exe audit CASE-2026-AETNA-01

# Output formal statutory appeal brief for a precedent case
./overrule.exe appeal CASE-2026-AETNA-01 -o elena_rostova_appeal.md
```

#### 4. Launch the Local Air-Gapped Web Dashboard

```bash
./overrule.exe serve --port 5050
```

Open `http://127.0.0.1:5050` in your browser to interact with the local visual interface.

---

## CLI Command Reference

| Command | Description |
|---|---|
| `overrule parse <file>` | Parse an EOB denial notice or ANSI 835 EDI remittance file |
| `overrule generate <file> [-o brief.md]` | Parse a denial notice and generate a complete court-ready appeal packet |
| `overrule list` | List all preloaded clinical demonstration cases |
| `overrule audit <case-id>` | Inspect carrier algorithmic flaws and clinical evidence citations |
| `overrule appeal <case-id> [-o brief.md]` | Compile a court-ready brief for a demonstration case |
| `overrule serve [--port 5050]` | Launch the local web console and REST API |
| `overrule help` | Display help and usage information |

---

## Repository Structure

```
overrule/
├── cmd/
│   ├── overrule/          # Main CLI application entrypoint
│   └── wasm/              # WebAssembly bridge for browser execution
├── examples/
│   ├── eob_sample.txt     # Real sanitized Aetna Proton Therapy denial letter
│   └── edi835_sample.edi  # Real ANSI X12 835 electronic remittance advice
├── pkg/
│   ├── appeal/            # Statutory ERISA appeal packet compiler
│   ├── clinical/          # Clinical guideline parity & PubMed evidence index
│   ├── denial/            # Clinical case schema and precedent records
│   ├── legal/             # Federal ERISA § 503 statutory rule definitions
│   ├── parser/            # EOB text regex & ANSI X12 835 EDI parsers
│   ├── rules/             # Carrier algorithmic vulnerability catalog (Optum, Cigna)
│   └── server/            # Embedded local-first REST API server
├── .github/workflows/     # Automated GitHub Actions CI workflow
├── LICENSE                # GNU Affero General Public License v3.0 (AGPL-3.0)
└── README.md
```

---

## Zero-PHI Privacy Guarantee

OVERRULE is built on a sovereign, local-first architecture:
- **Zero Cloud Transmission**: All OCR parsing, CARC/RARC code extraction, and legal brief assembly execute entirely in local process memory.
- **No Third-Party Telemetry**: No tracking pixels, third-party analytics, or remote API calls.
- **Air-Gapped Operation**: You can disconnect your network connection completely and the engine remains 100% operational.

---

## Legal & Statutory Foundation

This engine automates administrative remedies pursuant to:
- **ERISA § 503 (29 U.S.C. § 1133)**: Mandating adequate written notice and full, fair review by an appropriate named fiduciary.
- **29 C.F.R. § 2560.503-1(h)(2)(iii)**: Requiring plan administrators to provide, upon request and free of charge, all documents, records, and information relevant to the claim.
- **29 C.F.R. § 2560.503-1(h)(3)(iii)**: Requiring that reviews of medical necessity denials consult with a healthcare professional who has appropriate training and experience in the specific field of medicine.
- **ACA § 2719 (42 U.S.C. § 300gg-19)**: Guaranteeing independent external review parity for adverse benefit determinations.

---

## License

OVERRULE is licensed under the **[GNU Affero General Public License v3.0](LICENSE)** (AGPL-3.0).

This copyleft license guarantees that the core patient defense engine remains free and open source forever. Any entity incorporating or deploying this engine over a network must make their complete source code publicly available under the same license terms.
