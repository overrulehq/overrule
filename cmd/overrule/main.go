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
	"github.com/overrulehq/overrule/pkg/server"
)

const banner = `
   ____  _   _ _____ ____  ____  _   _ _     _____ 
  / __ \| | | | ____|  _ \|  _ \| | | | |   | ____|
 | |  | | | | |  _| | |_) | |_) | | | | |   |  _|  
 | |__| | |_| | |___|  _ <|  _ <| |_| | |___| |___ 
  \____/ \___/|_____|_| \_\_| \_\\___/|_____|_____|
  The Autonomous Patient Defense & Denial Overturn Engine
  Built under Apache 2.0 • Local-First Privacy Guarantee
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
	fmt.Println("  serve [--port 5050]        Launch interactive patient defense web console")
	fmt.Println("  list                       List built-in clinical denial demonstration cases")
	fmt.Println("  audit <case-id>            Run clinical and statutory vulnerability audit")
	fmt.Println("  appeal <case-id> [-o path] Generate complete court-ready statutory appeal brief")
	fmt.Println("  help                       Show this help message")
}
