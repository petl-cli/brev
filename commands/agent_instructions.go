package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// agentInstructionsCmd prints the llms.txt content at runtime so agents can
// include it in their system prompt without needing a separate file:
//
//	INSTRUCTIONS=$(brevo-api agent-instructions)
var agentInstructionsCmd = &cobra.Command{
	Use:   "agent-instructions",
	Short: "Print machine-readable instructions for AI agents (llms.txt format)",
	Long: `Prints a complete description of this CLI's commands, flags, exit codes,
and usage patterns optimised for inclusion in an AI agent's system prompt.

Example:
  # Include in Claude Code context:
  brevo-api agent-instructions > CLAUDE.md

  # Capture inline:
  INSTRUCTIONS=$(brevo-api agent-instructions)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(agentInstructionsContent)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentInstructionsCmd)
}

// agentInstructionsContent is the full llms.txt baked into the binary at build time.
// Regenerate by re-running the CLI generator against the same OpenAPI spec.
const agentInstructionsContent = `# brevo-api

Brevo provide a RESTFul API that can be used with any languages. With this API, you will be able to :
  - Manage your campaigns and get the statistics
  - Manage your contacts
  - Send transactional Emails and SMS
  - and much more...

You can download our wrappers at https://github.com/orgs/brevo

**Possible responses**
  | Code | Message |
  | :-------------: | ------------- |
  | 200  | OK. Successful Request  |
  | 201  | OK. Successful Creation |
  | 202  | OK. Request accepted |
  | 204  | OK. Successful Update/Deletion  |
  | 400  | Error. Bad Request  |
  | 401  | Error. Authentication Needed  |
  | 402  | Error. Not enough credit, plan upgrade needed  |
  | 403  | Error. Permission denied  |
  | 404  | Error. Object does not exist |
  | 405  | Error. Method not allowed  |
  | 406  | Error. Not Acceptable  |


This file is the agent-facing overview of the ` + "`" + `brevo-api` + "`" + ` CLI. It explains what the tool does and how to use it well — *not* every flag of every command. For per-command details run ` + "`" + `brevo-api <command> --help` + "`" + ` or ` + "`" + `brevo-api <command> --schema` + "`" + ` (returns JSON).

## Install

The binary is ` + "`" + `brevo-api` + "`" + `. Build from source or download a release.

## Authentication

- API key — set ` + "`" + `BREVO_API_API_KEY` + "`" + ` or pass ` + "`" + `--api-key <key>` + "`" + `

Run ` + "`" + `brevo-api configure` + "`" + ` to set credentials interactively.

## Commands

brevo-api groups operations by resource. One or two examples per group below — for the full set of commands and flags use ` + "`" + `brevo-api <group> --help` + "`" + `.

### ` + "`" + `account` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api account get-user-activity-logs    # Get user activity logs
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api account information-details    # Get your account information, plan and credits details
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api account --help` + "`" + `

### ` + "`" + `companies` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api companies create-company --body '{...}'    # Create a company
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api companies delete-company --id <string>    # Delete a company
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api companies --help` + "`" + `

### ` + "`" + `contacts` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api contacts add-to-list --list-id <integer> --body '{...}'    # Add existing contacts to a list
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api contacts create-attribute --attribute-category <string> --attribute-name <string> --body '{...}'    # Create contact attribute
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api contacts --help` + "`" + `

### ` + "`" + `conversations` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api conversations delete-automated-message --id <string>    # Delete an automated message
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api conversations delete-message-sent-by-agent --id <string>    # Delete a message sent by an agent
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api conversations --help` + "`" + `

### ` + "`" + `coupons` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api coupons create-collection --body '{...}'    # Create а coupon collection
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api coupons create-coupon-collection --body '{...}'    # Create coupons for a coupon collection
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api coupons --help` + "`" + `

### ` + "`" + `deals` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api deals create-new-deal --body '{...}'    # Create a deal
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api deals delete-deal --id <string>    # Delete a deal
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api deals --help` + "`" + `

### ` + "`" + `domains` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api domains authenticate-domain --domain-name <string>    # Authenticate a domain
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api domains create-new-domain    # Create a new domain
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api domains --help` + "`" + `

### ` + "`" + `ecommerce` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api ecommerce activate-app    # Activate the eCommerce app
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api ecommerce create-categories-batch --body '{...}'    # Create categories in batch
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api ecommerce --help` + "`" + `

### ` + "`" + `email-campaigns` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api email-campaigns create-campaign --body '{...}'    # Create an email campaign
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api email-campaigns export-recipients-post --campaign-id <integer>    # Export the recipients of an email campaign
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api email-campaigns --help` + "`" + `

### ` + "`" + `event` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api event track-interaction --body '{...}'    # Create an event
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api event --help` + "`" + `

### ` + "`" + `external-feeds` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api external-feeds create-feed --body '{...}'    # Create an external feed
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api external-feeds delete-feed-by-uuid --uuid <string>    # Delete an external feed
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api external-feeds --help` + "`" + `

### ` + "`" + `files` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api files delete-file --id <string>    # Delete a file
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api files download-file --id <string>    # Download a file
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api files --help` + "`" + `

### ` + "`" + `inbound-parsing` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api inbound-parsing get-all-events    # Get the list of all the events for the received emails.
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api inbound-parsing get-attachment-by-token --download-token <string>    # Retrieve inbound attachment with download token.
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api inbound-parsing --help` + "`" + `

### ` + "`" + `master-account` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api master-account account-get-user-activity-logs    # Get user activity logs
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api master-account check-admin-user-permissions --email <string>    # Check admin user permissions
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api master-account --help` + "`" + `

### ` + "`" + `notes` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api notes create-new-note --body '{...}'    # Create a note
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api notes get-all    # Get all notes
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api notes --help` + "`" + `

### ` + "`" + `payments` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api payments create-payment-request --body '{...}'    # Create a payment request
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api payments delete-request --id <string>    # Delete a payment request.
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api payments --help` + "`" + `

### ` + "`" + `process` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api process get-all-processes    # Return all the processes for your account
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api process get-process-information --process-id <integer>    # Return the informations for a process
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api process --help` + "`" + `

### ` + "`" + `reseller` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api reseller add-child-credits --child-identifier <string> --body '{...}'    # Add Email and/or SMS credits to a specific child account
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api reseller associate-dedicated-ip-to-child --child-identifier <string> --body '{...}'    # Associate a dedicated IP to the child
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api reseller --help` + "`" + `

### ` + "`" + `senders` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api senders create-new-sender    # Create a new sender
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api senders get-all-dedicated-ips --sender-id <integer>    # Get all the dedicated IPs for a sender
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api senders --help` + "`" + `

### ` + "`" + `sms-campaigns` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api sms-campaigns create-campaign --body '{...}'    # Creates an SMS campaign
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api sms-campaigns export-recipients-process --campaign-id <integer>    # Export an SMS campaign's recipients
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api sms-campaigns --help` + "`" + `

### ` + "`" + `tasks` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api tasks create-new-task --body '{...}'    # Create a task
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api tasks get-all    # Get all tasks
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api tasks --help` + "`" + `

### ` + "`" + `transactional-emails` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api transactional-emails add-blocked-domain --body '{...}'    # Add a new domain to the list of blocked domains
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api transactional-emails create-template --body '{...}'    # Create an email template
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api transactional-emails --help` + "`" + `

### ` + "`" + `transactional-sms` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api transactional-sms get-aggregated-report    # Get your SMS activity aggregated over a period of time
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api transactional-sms get-all-events    # Get all your SMS activity (unaggregated events)
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api transactional-sms --help` + "`" + `

### ` + "`" + `transactional-whatsapp` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api transactional-whatsapp get-activity    # Get all your WhatsApp activity (unaggregated events)
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api transactional-whatsapp send-message --body '{...}'    # Send a WhatsApp message
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api transactional-whatsapp --help` + "`" + `

### ` + "`" + `user` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api user check-permission --email <string>    # Check user permission
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api user get-all-users    # Get the list of all your users
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api user --help` + "`" + `

### ` + "`" + `webhooks` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api webhooks create-hook --body '{...}'    # Create a webhook
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api webhooks delete-webhook --webhook-id <integer>    # Delete a webhook
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api webhooks --help` + "`" + `

### ` + "`" + `whatsapp-campaigns` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api whatsapp-campaigns create-and-send --body '{...}'    # Create and Send a WhatsApp campaign
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api whatsapp-campaigns create-template --body '{...}'    # Create a WhatsApp template
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `brevo-api whatsapp-campaigns --help` + "`" + `


## Output and parsing

All output is JSON by default and goes to stdout. Errors go to stderr as a single-line JSON object — stdout stays clean for piping.

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api <cmd> --jq <path>          # extract fields without jq installed (GJSON syntax)
brevo-api <cmd> -o yaml              # change format: json (default), yaml, table, compact, raw, pretty
brevo-api <cmd> --dry-run            # print the HTTP request without sending it
brevo-api <cmd> --schema             # JSON schema of inputs and outputs for that command
` + "`" + `` + "`" + `` + "`" + `

GJSON path examples:
- ` + "`" + `--jq id` + "`" + ` — scalar
- ` + "`" + `--jq items.#.id` + "`" + ` — every id from an array
- ` + "`" + `--jq "items.#(active==true)#"` + "`" + ` — filter array by condition

## Exit codes

Branch on ` + "`" + `$?` + "`" + ` rather than parsing stderr:

| Code | Meaning |
|------|---------|
| 0 | success |
| 1 | unknown error |
| 2 | auth failed (401 / 403) |
| 3 | not found (404) |
| 4 | validation error (400 / 422) |
| 5 | rate limited (429) |
| 6 | server error (5xx) |
| 7 | network error |

Error JSON shape on stderr:
` + "`" + `` + "`" + `` + "`" + `json
{"error":true,"code":"not_found","status":404,"message":"...","exit_code":3}
` + "`" + `` + "`" + `` + "`" + `

## Common workflows

### List then fetch one
` + "`" + `` + "`" + `` + "`" + `bash
ID=$(brevo-api <group> list --jq "items.0.id")
brevo-api <group> get --id "$ID"
` + "`" + `` + "`" + `` + "`" + `

### Capture a created ID
` + "`" + `` + "`" + `` + "`" + `bash
ID=$(brevo-api <group> create [flags] --jq id)
` + "`" + `` + "`" + `` + "`" + `

### Safe destructive call
` + "`" + `` + "`" + `` + "`" + `bash
brevo-api <group> delete --id X --dry-run    # inspect first
brevo-api <group> delete --id X              # exit 0 = deleted, 3 = not found
` + "`" + `` + "`" + `` + "`" + `

### Branch on exit code
` + "`" + `` + "`" + `` + "`" + `bash
brevo-api <cmd> [flags]
case $? in
  0) : success ;;
  2) echo "fix credentials" ;;
  3) echo "not found, skip" ;;
  5) sleep 10 ; retry ;;
esac
` + "`" + `` + "`" + `` + "`" + `

## Discovering more

` + "`" + `` + "`" + `` + "`" + `bash
brevo-api --help                              # top-level groups
brevo-api <group> --help                      # commands in a group
brevo-api <group> <cmd> --help                # flags + description
brevo-api <group> <cmd> --schema              # JSON schema
brevo-api agent-instructions                  # this file, embedded in the binary
` + "`" + `` + "`" + `` + "`" + `

`
