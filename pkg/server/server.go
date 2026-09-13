package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"html"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/overrulehq/overrule/pkg/appeal"
	"github.com/overrulehq/overrule/pkg/denial"
	"github.com/overrulehq/overrule/pkg/rules"
)

//go:embed public/*
var publicFS embed.FS

// Server manages the HTTP REST endpoints and embedded frontend.
type Server struct {
	mu    sync.RWMutex
	cases map[string]*denial.DenialCase
	port  int
}

// NewServer initializes a server preloaded with clinical test cases.
func NewServer(port int) *Server {
	s := &Server{
		cases: make(map[string]*denial.DenialCase),
		port:  port,
	}

	// Preload sample cases
	for i := range denial.SampleCases {
		c := denial.SampleCases[i]
		s.cases[c.ID] = &c
	}

	return s
}

// Start boots the HTTP server.
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/cases", s.handleCases)
	mux.HandleFunc("/api/cases/", s.handleCaseByID)
	mux.HandleFunc("/api/audit", s.handleAudit)
	mux.HandleFunc("/api/audit/", s.handleAudit)
	mux.HandleFunc("/api/appeal", s.handleAppeal)
	mux.HandleFunc("/api/appeal/", s.handleAppeal)
	mux.HandleFunc("/api/payers", s.handlePayers)

	// Static frontend
	subFS, err := fs.Sub(publicFS, "public")
	if err != nil {
		return fmt.Errorf("failed to load public sub fs: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(subFS)))

	addr := fmt.Sprintf("127.0.0.1:%d", s.port)
	log.Printf("[OVERRULE] Patient Defense Engine active at http://%s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodGet {
		s.mu.RLock()
		defer s.mu.RUnlock()
		list := make([]*denial.DenialCase, 0, len(s.cases))
		for _, c := range s.cases {
			list = append(list, c)
		}
		json.NewEncoder(w).Encode(list)
		return
	}

	if r.Method == http.MethodPost {
		var newCase denial.DenialCase
		if err := json.NewDecoder(r.Body).Decode(&newCase); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}
		if newCase.ID == "" {
			newCase.ID = fmt.Sprintf("CASE-%d", len(s.cases)+1)
		}
		s.mu.Lock()
		s.cases[newCase.ID] = &newCase
		s.mu.Unlock()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCase)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleCaseByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := strings.TrimPrefix(r.URL.Path, "/api/cases/")
	s.mu.RLock()
	c, exists := s.cases[id]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Case not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(c)
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := r.URL.Query().Get("id")
	if id == "" {
		id = strings.TrimPrefix(r.URL.Path, "/api/audit/")
		id = strings.Trim(id, "/")
	}

	s.mu.RLock()
	c, exists := s.cases[id]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Case not found", http.StatusNotFound)
		return
	}

	packet := appeal.GeneratePacket(c)

	// Format vulnerabilities for API
	type VulnDTO struct {
		Carrier        string `json:"carrier"`
		AlgorithmName  string `json:"algorithm_name"`
		FlawType       string `json:"flaw_type"`
		Description    string `json:"description"`
		LegalGround    string `json:"legal_ground"`
		CasePrecedent  string `json:"case_precedent"`
		ErisaDemand    string `json:"erisa_demand"`
	}

	vulns := make([]VulnDTO, 0, len(packet.Vulnerabilities))
	for _, v := range packet.Vulnerabilities {
		vulns = append(vulns, VulnDTO{
			Carrier:       v.PayerName,
			AlgorithmName: v.RuleName,
			FlawType:      v.Severity,
			Description:   v.Description,
			LegalGround:   v.StatutoryBasis,
			CasePrecedent: v.ReversalPrecedent,
			ErisaDemand:   v.RecommendedAction,
		})
	}

	// Clinical evidence
	type PubDTO struct {
		PMID       string `json:"pmid"`
		Title      string `json:"title"`
		Authors    string `json:"authors"`
		Journal    string `json:"journal"`
		Year       int    `json:"year"`
		Conclusion string `json:"conclusion"`
	}

	pubs := make([]PubDTO, 0)
	guidelines := make([]string, 0)
	standardOfCare := "Established specialty society guideline parity."

	for _, ca := range packet.ClinicalArguments {
		if ca.StandardOfCare != "" {
			standardOfCare = ca.StandardOfCare
		}
		if ca.ClinicalRationale != "" {
			guidelines = append(guidelines, ca.ClinicalRationale)
		}
		for _, lit := range ca.Literature {
			pubs = append(pubs, PubDTO{
				PMID:       lit.PMID,
				Title:      lit.Citation,
				Authors:    lit.Grade,
				Journal:    lit.Journal,
				Year:       2026,
				Conclusion: lit.KeyFinding,
			})
		}
	}

	resp := map[string]interface{}{
		"case":            c,
		"vulnerabilities": vulns,
		"clinical_evidence": map[string]interface{}{
			"medical_guidelines":    guidelines,
			"pubmed_citations":      pubs,
			"standard_of_care":      standardOfCare,
			"physician_declaration": packet.PhysicianRebuttal,
		},
		"recommended_statute": "ERISA § 503 (29 U.S.C. § 1133)",
	}

	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleAppeal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := r.URL.Query().Get("id")
	if id == "" {
		path := strings.TrimPrefix(r.URL.Path, "/api/appeal/")
		parts := strings.Split(path, "/")
		id = parts[0]
	}

	s.mu.RLock()
	c, exists := s.cases[id]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Case not found", http.StatusNotFound)
		return
	}

	packet := appeal.GeneratePacket(c)
	md := packet.ToMarkdown()

	// Check if raw markdown export requested
	if strings.Contains(r.URL.Path, "/markdown") || r.URL.Query().Get("format") == "markdown" {
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"Overrule-Appeal-%s.md\"", c.ID))
		w.Write([]byte(md))
		return
	}

	// Simple rich HTML representation of the markdown
	htmlContent := markdownToHTML(md)

	resp := map[string]interface{}{
		"id":               c.ID,
		"patient_name":     c.PatientName,
		"payer_name":       c.InsurerName,
		"statutory_grounds": []string{"29 U.S.C. § 1133", "29 C.F.R. § 2560.503-1", "ACA § 2719"},
		"appeal_markdown":  md,
		"appeal_html":      htmlContent,
	}

	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handlePayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(rules.Catalog)
}

func markdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var out strings.Builder
	inCode := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			if !inCode {
				inCode = true
				out.WriteString("<pre style=\"background: var(--wabi-paper-alt); border: 1px solid var(--border-organic); padding: 1.25rem; border-radius: 8px; font-family: var(--font-mono); font-size: 0.85rem; overflow-x: auto; margin: 1.25rem 0; white-space: pre-wrap;\"><code>")
			} else {
				inCode = false
				out.WriteString("</code></pre>\n")
			}
			continue
		}

		if inCode {
			out.WriteString(html.EscapeString(line))
			out.WriteString("\n")
			continue
		}

		if strings.HasPrefix(trimmed, "# ") {
			out.WriteString(fmt.Sprintf("<h1 style=\"font-size: 2rem; font-family: var(--font-serif); margin: 2rem 0 1rem; color: var(--ink-black);\">%s</h1>\n", html.EscapeString(trimmed[2:])))
		} else if strings.HasPrefix(trimmed, "## ") {
			out.WriteString(fmt.Sprintf("<h2 style=\"font-size: 1.5rem; font-family: var(--font-serif); margin: 1.75rem 0 0.75rem; border-bottom: 1px solid var(--border-organic); padding-bottom: 0.4rem; color: var(--ink-black);\">%s</h2>\n", html.EscapeString(trimmed[3:])))
		} else if strings.HasPrefix(trimmed, "### ") {
			out.WriteString(fmt.Sprintf("<h3 style=\"font-size: 1.2rem; font-family: var(--font-serif); margin: 1.25rem 0 0.5rem; color: var(--ink-black);\">%s</h3>\n", html.EscapeString(trimmed[4:])))
		} else if strings.HasPrefix(trimmed, "> ") {
			out.WriteString(fmt.Sprintf("<blockquote style=\"border-left: 3px solid var(--clay-persimmon); padding-left: 1rem; color: var(--ink-muted); margin: 1rem 0; font-style: italic;\">%s</blockquote>\n", html.EscapeString(trimmed[2:])))
		} else if strings.HasPrefix(trimmed, "- ") {
			out.WriteString(fmt.Sprintf("<li style=\"margin-left: 1.5rem; margin-bottom: 0.4rem; color: var(--ink-muted);\">%s</li>\n", html.EscapeString(trimmed[2:])))
		} else if trimmed == "---" {
			out.WriteString("<hr style=\"border: none; border-top: 1px solid var(--border-organic); margin: 2rem 0;\" />\n")
		} else if trimmed != "" {
			out.WriteString(fmt.Sprintf("<p style=\"margin-bottom: 1rem; color: var(--ink-black); line-height: 1.75;\">%s</p>\n", html.EscapeString(trimmed)))
		}
	}

	return out.String()
}
