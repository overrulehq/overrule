# OVERRULE Technical & Product Roadmap (2026–2028)
### The Autonomous Patient Defense & Algorithmic Claim Reversal Engine

---

## 1. Executive Summary & Core Mission
**OVERRULE** (github.com/overrulehq/overrule) is the world's first open-source, local-first legal and clinical engine designed to deconstruct algorithmic health insurance claim rejections and generate court-ready administrative appeal dossiers under **ERISA § 503 (29 U.S.C. § 1133)** and the **Affordable Care Act § 2719**.

Health insurance conglomerates reject over 200 million claims annually using automated batch models (**Optum nH Predict**, **Cigna PXDX**, **Aetna CPB Automated Triage**) in 1.2 seconds without meaningful physician review. Fewer than 0.2% of patients appeal due to administrative exhaustion.

OVERRULE is built to shift the economic equation: **Make automated denial an untenable financial and regulatory liability for insurance carriers by arming every patient with an automated federal legal weapon.**

---

## 2. Competitive Landscape & Defensibility Moat

### Why There Is No Direct Competitor to OVERRULE
| Capability | Generic AI Appeal Tools (e.g. Claimable, GPT wrappers) | Commercial Revenue Cycle Tools (Change, Waystar) | **OVERRULE (Open-Source Core)** |
| :--- | :--- | :--- | :--- |
| **Legal Strategy** | Polite conversational 'pleas' (ignored by carriers) | Hospital billing maximization (serves payers/providers) | **Aggressive Federal ERISA § 503 Discovery Demands** |
| **Privacy / PHI** | Transmits sensitive health data to third-party clouds | Enterprise BAA cloud silos | **100% Air-Gapped Local-First Privacy (Zero PHI Egress)** |
| **Carrier Dissection** | None (treats denials as generic errors) | Proprietary clearinghouse rules | **Dedicated Algorithmic Flaw Catalog (UHC, Aetna, Cigna, BCBS)** |
| **Medical Parity** | Hallucinated or broad medical summaries | Billing code crosswalks | **Binding PubMed PMIDs & Clinical Specialty Society Standards** |
| **Distribution & Trust**| Closed-source paid paywalls | B2B enterprise sales cycles | **Free Sovereign Open-Source Core (Community Verified)** |

---

## 3. Licensing Strategy: Why AGPLv3 vs. Apache 2.0

### The Recommendation: **GNU Affero General Public License v3 (AGPLv3)**
- **Why AGPLv3?**
  In healthcare tech, predatory venture-funded startups or health insurance conglomerates themselves often take open-source tools, wrap them behind a proprietary paid paywall or cloud API, and commoditize the open-source creators without giving back.
- Under **AGPLv3**, anyone who runs OVERRULE—even as a network service or cloud API—**must release their complete source code and algorithmic rule additions back to the public**.
- This permanently protects the patient community: No insurer or predatory health tech intermediary can privatize the legal defenses of patients.
- Developers and non-profits can freely inspect, run, contribute, and extend the engine locally.

*(Alternative: Apache 2.0 if the strategic goal is maximum enterprise developer adoption where commercial employers embed the library into internal benefits portals without copyleft obligations).*

---

## 4. The Steve Jobs Patient UX Philosophy ('Simplicity as Sanctuary')

If Steve Jobs examined the state of healthcare appeals, his diagnosis would be immediate:
> *'A frightened patient with a cancer diagnosis or a parent with a hospitalized child doesn't want a legal dashboard, a database grid, or 10 CARC codes. They have one question: " Is my treatment covered and what do I do next?\ Everything else is noise.'*

### The 3-Second Patient Interaction Loop:
1. **The Single Canvas**: One clear focal point. Drag-and-drop the denial letter or snap a phone photo.
2. **Instant Radical Clarity**: Within 3 seconds, local OCR displays:
 - *'Aetna denied your Brain MRI citing an automated rule.'*
 - *'Under Federal ERISA law, this denial is invalid. Reversal likelihood: 96.4%.'*
3. **The Single Action Button**: *'Overrule This Denial'*.
 - Generates the certified 12-page legal folio, attaches physician declarations, and prepares the physical mailing packet with one tap.

---

## 5. Development Phases

### Phase 1: Sovereign Go Engine & Statutory Core (Completed)
- [x] High-performance Go 1.26 CLI binary (overrule).
- [x] Algorithmic vulnerability catalogs for UnitedHealthcare, Aetna, Cigna, and BlueCross.
- [x] Comprehensive 12-page ERISA § 503 & ACA § 2719 court-ready legal brief builder.
- [x] Built-in REST API and embedded web serving on port 5050.
- [x] Zero-cloud local-first privacy guarantee.

### Phase 2: On-Device Vision & Clinical Guideline Live Sync (Q4 2026)
- [ ] **WASM / Local Tesseract OCR Pipeline**: Instant camera photo and scanned PDF parsing of EOBs directly on the device.
- [ ] **ICD-10 / CPT Automated Tokenizer**: Auto-extracts diagnosis and procedure codes from messy doctor notes.
- [ ] **Live PubMed & OpenAlex Parity Fetcher**: Programmatically retrieves latest Level-1A clinical trial abstracts from NIH PubMed matching patient conditions.
- [ ] **Doctor Magic Link**: A secure, 1-click web signature sheet sent to the treating clinic so the physician can sign the pre-written clinical necessity declaration in 10 seconds.

### Phase 3: Programmatic Transmission & Regulatory Cross-Filing (Q1 2027)
- [ ] **USPS Certified Mail Programmatic API (Lob / PostGrid)**: 
 Automatically print and ship certified physical appeal packets with USPS Form 3800 tracking directly from the console.
- [ ] **State Insurance Commissioner (DOI) Auto-Filer**:
 Simultaneously cross-file formal consumer complaints with state regulators (e.g., California CDI, New York DFS, Texas TDI) to trigger mandatory 15-day state review deadlines.
- [ ] **ERISA § 502(a)(1)(B) Federal Court Civil Complaint Generator**:
 If internal administrative appeals are exhausted, format a court-ready federal complaint for filing in U.S. District Court.

### Phase 4: B2B Enterprise & Self-Funded Employer Fiduciary Shield (Q2 2027 — YC Commercialization)
- [ ] **The Problem**: 65% of covered US workers are in self-funded plans where the *employer* pays the claims and bears personal legal fiduciary liability under ERISA (*Kraft Heinz v. Aetna*).
- [ ] **The Solution**: **OVERRULE Sentinel** — Enterprise B2B SaaS for Chief Human Resources Officers & Plan Fiduciaries:
 - Continuously audits TPA claim rejections across the company's workforce.
 - Automatically reverses improper denials before employees face collection agencies.
 - Price: .50 – .00 PEPM (Per Employee Per Month).

### Phase 5: National Algorithmic Audit Registry (2028)
- [ ] Aggregation of zero-knowledge, privacy-preserving procedural denial patterns.
- [ ] Public transparency index documenting insurer algorithmic error rates to aid Department of Justice antitrust investigations and congressional oversight.

---

## 6. Open-Source Community Guidelines
- **Contributions**: Pull requests for new payer profiles (ules/payers/*.go) and clinical trial mappings are welcome.
- **Security**: Because OVERRULE handles medical dispute records, all core cryptographic and parsing routines must maintain zero external telemetry.
