package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/overrulehq/overrule/pkg/appeal"
	"github.com/overrulehq/overrule/pkg/denial"
	"github.com/overrulehq/overrule/pkg/parser"
	"github.com/overrulehq/overrule/pkg/server"
)

const banner = `
   ____  _   _ _____ ____  ____  _   _ _     _____ 
  / __ \| | | | ____|  _ \|  _ \| | | | |   | ____|
 | |  | | | | |  _| | |_) | |_) | | | | |   |  _|  
 | |__| | |_| | |___|  _ <|  _ <| |_| | |___| |___ 
  \____/ \___/|_____|_| \_\_| \_\\___/|_____|_____|
  The Autonomous Patient Defense & Denial Overturn Engine
  Licensed under AGPL-3.0 • Local-First Privacy Guarantee
`

func main() {
	if len(os.Args) < 2 {
		runServer(5050, true)
		return
	}

	command := os.Args[1]

	switch command {
	case "serve", "server", "web":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		port := serveCmd.Int("port", 5050, "Port for web dashboard")
		noBrowser := serveCmd.Bool("no-browser", false, "Do not auto-launch browser")
		serveCmd.Parse(os.Args[2:])
		runServer(*port, !*noBrowser)

	case "list":
		fmt.Print(banner)
		fmt.Println("Available Demonstration Cases:")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, c := range denial.SampleCases {
			fmt.Printf("[%s] %s | Payer: %-18s | Denied: $%-9.2f | Dx: %s\n",
				c.ID, c.PatientName, c.InsurerName, c.DeniedAmount, c.Diagnoses[0].Code)
		}
		fmt.Println("--------------------------------------------------------------------------------")

	case "parse":
		if len(os.Args) < 3 {
			fmt.Println("Usage: overrule parse <eob.txt | remittance.edi>")
			os.Exit(1)
		}
		filePath := os.Args[2]
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		content := string(data)
		if strings.HasSuffix(strings.ToLower(filePath), ".edi") || strings.HasPrefix(strings.TrimSpace(content), "ISA*") {
			claims, err := parser.Parse835EDI(content)
			if err != nil {
				fmt.Printf("Error parsing 835 EDI: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully parsed %d claim denials from 835 Remittance File:\n\n", len(claims))
			for i, claim := range claims {
				c := claim.ToDenialCase()
				fmt.Printf("[%d] Claim: %s | Patient: %s | Payer: %s | Denied: $%.2f (CARC: %s)\n",
					i+1, c.ClaimNumber, c.PatientName, c.InsurerName, c.DeniedAmount, c.DenialReason.CARCCode)
			}
		} else {
			parsed := parser.ParseRawText(content)
			c := parsed.ToDenialCase("Elena Rostova", "Dr. Catherine M. Lee, MD")
			fmt.Printf("=== OVERRULE PARSER VERDICT ===\n")
			fmt.Printf("Payer:         %s\n", c.InsurerName)
			fmt.Printf("Patient:       %s\n", c.PatientName)
			fmt.Printf("Claim Number:  %s\n", c.ClaimNumber)
			fmt.Printf("Policy ID:     %s\n", c.PolicyID)
			fmt.Printf("Denied Amount: $%.2f\n", c.DeniedAmount)
			fmt.Printf("Denial Code:   %s (%s)\n", c.DenialReason.CARCCode, c.DenialReason.Category)
			fmt.Printf("Doctor:        %s (%s)\n", c.Physician.Name, c.Physician.Specialty)
		}

	case "generate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: overrule generate <file.txt | file.edi> [-o output.md]")
			os.Exit(1)
		}
		filePath := os.Args[2]
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		content := string(data)
		var c *denial.DenialCase

		if strings.HasSuffix(strings.ToLower(filePath), ".edi") || strings.HasPrefix(strings.TrimSpace(content), "ISA*") {
			claims, err := parser.Parse835EDI(content)
			if err != nil || len(claims) == 0 {
				fmt.Printf("Error parsing 835 EDI file: %v\n", err)
				os.Exit(1)
			}
			c = claims[0].ToDenialCase()
		} else {
			parsed := parser.ParseRawText(content)
			c = parsed.ToDenialCase("Elena Rostova", "Dr. Catherine M. Lee, MD")
		}

		packet := appeal.GeneratePacket(c)
		md := packet.ToMarkdown()

		outFile := ""
		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "-o" && i+1 < len(os.Args) {
				outFile = os.Args[i+1]
				break
			}
		}

		if outFile != "" {
			if err := os.WriteFile(outFile, []byte(md), 0644); err != nil {
				fmt.Printf("Error writing file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully generated court-ready appeal brief for %s ($%.2f) to: %s\n",
				c.PatientName, c.DeniedAmount, outFile)
		} else {
			fmt.Println(md)
		}

	case "audit":
		if len(os.Args) < 3 {
			fmt.Println("Usage: overrule audit <CASE-ID>")
			os.Exit(1)
		}
		caseID := os.Args[2]
		c := findCase(caseID)
		if c == nil {
			fmt.Printf("Error: Case %s not found. Run 'overrule list' to see available cases.\n", caseID)
			os.Exit(1)
		}

		packet := appeal.GeneratePacket(c)
		fmt.Print(banner)
		fmt.Printf("=== OVERRULE AUDIT: %s (%s) ===\n", c.PatientName, c.ID)
		fmt.Printf("Payer: %s | Claim: %s | Amount Denied: $%.2f\n", c.InsurerName, c.ClaimNumber, c.DeniedAmount)
		fmt.Printf("Alleged Denial Reason: %s\n\n", c.DenialReason.RawLetterText)

		fmt.Println("--- IDENTIFIED ALGORITHMIC & STATUTORY VULNERABILITIES ---")
		for _, v := range packet.Vulnerabilities {
			fmt.Printf("[%s] %s\n", v.Severity, v.RuleName)
			fmt.Printf("  Finding:   %s\n", v.Description)
			fmt.Printf("  Statute:   %s\n", v.StatutoryBasis)
			fmt.Printf("  Precedent: %s\n\n", v.ReversalPrecedent)
		}

		fmt.Println("--- CLINICAL EVIDENCE & GUIDELINE PARITY ---")
		for _, ca := range packet.ClinicalArguments {
			fmt.Printf("Target: %s for %s\n", ca.ProcedureCode, ca.DiagnosisCode)
			fmt.Printf("  Standard: %s\n", ca.StandardOfCare)
			for _, lit := range ca.Literature {
				fmt.Printf("  * %s (%s, PMID: %s)\n", lit.Citation, lit.Journal, lit.PMID)
			}
		}

	case "appeal":
		if len(os.Args) < 3 {
			fmt.Println("Usage: overrule appeal <CASE-ID> [-o output.md]")
			os.Exit(1)
		}
		caseID := os.Args[2]
		c := findCase(caseID)
		if c == nil {
			fmt.Printf("Error: Case %s not found. Run 'overrule list' to see available cases.\n", caseID)
			os.Exit(1)
		}

		packet := appeal.GeneratePacket(c)
		md := packet.ToMarkdown()

		if len(os.Args) >= 5 && os.Args[3] == "-o" {
			outFile := os.Args[4]
			if err := os.WriteFile(outFile, []byte(md), 0644); err != nil {
				fmt.Printf("Error writing file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully generated court-ready appeal brief to: %s\n", outFile)
		} else {
			fmt.Println(md)
		}

	case "help", "--help", "-h":
		printHelp()

	default:
		fmt.Printf("Unknown command '%s'. Run 'overrule help' for usage.\n", command)
		os.Exit(1)
	}
}

func findCase(id string) *denial.DenialCase {
	for i := range denial.SampleCases {
		if strings.EqualFold(denial.SampleCases[i].ID, id) {
			return &denial.SampleCases[i]
		}
	}
	return nil
}

func runServer(port int, openBrowser bool) {
	fmt.Print(banner)
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	fmt.Printf("Starting Overrule Patient Defense Console on %s\n", url)
	fmt.Println("Press Ctrl+C to stop.")

	if openBrowser {
		go func() {
			switch runtime.GOOS {
			case "windows":
				exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
			case "darwin":
				exec.Command("open", url).Start()
			case "linux":
				exec.Command("xdg-open", url).Start()
			}
		}()
	}

	s := server.NewServer(port)
	if err := s.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(banner)
	fmt.Println("Usage: overrule <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  serve [--port 5050]                  Launch interactive patient defense web console")
	fmt.Println("  parse <file.txt | file.edi>          Parse an EOB denial notice or ANSI 835 EDI file")
	fmt.Println("  generate <file> [-o appeal.md]       Parse a file and generate a court-ready appeal brief")
	fmt.Println("  list                                 List built-in clinical denial demonstration cases")
	fmt.Println("  audit <case-id>                      Run clinical and statutory vulnerability audit")
	fmt.Println("  appeal <case-id> [-o path]           Generate brief for preloaded demonstration case")
	fmt.Println("  help                                 Show this help message")
}
