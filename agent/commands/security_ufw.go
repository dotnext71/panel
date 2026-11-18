package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/dotnext-panel/dotnext/agentwire"
	"github.com/dotnext-panel/dotnext/internal/platform/osinfo"
)

// M26 Step 2 (ADR-0054). UFW rule + status surface for the admin
// Security tab. Every handler validates inputs strictly BEFORE shelling
// out to ufw — the CLI is happy to accept "allow 22 OR DROP TABLE" and
// silently misinterpret it. The panel-api layer must trust nothing
// upstream.

// ---- shared helpers --------------------------------------------------------

func ufwInvalidArg(msg string) error {
	return &agentwire.AgentError{Code: agentwire.CodeInvalidArgument, Message: msg}
}

func ufwInternal(msg string, err error) error {
	return &agentwire.AgentError{
		Code:    agentwire.CodeInternal,
		Message: fmt.Sprintf("%s: %v", msg, err),
	}
}

// ---- security.ufw.status ---------------------------------------------------

type ufwRule struct {
	Num    int    `json:"num"`
	Action string `json:"action"`
	From   string `json:"from"`
	To     string `json:"to"`
	Proto  string `json:"proto,omitempty"`
	Port   string `json:"port,omitempty"`
}

type ufwStatusResponse struct {
	Active     bool      `json:"active"`
	DefaultIn  string    `json:"default_in"`
	DefaultOut string    `json:"default_out"`
	Rules      []ufwRule `json:"rules"`
}

func ufwStatusHandler(ctx context.Context, _ json.RawMessage) (any, error) {
	if isFirewalldBackend() {
		return firewalldStatusHandler(ctx)
	}
	// Two calls: `ufw status numbered verbose` drops the "Default:"
	// line that `ufw status verbose` (without numbered) prints. We need
	// the rule numbers (for delete-by-num) AND the default policies, so
	// query both.
	numbered, err := execCommandContext(ctx, "ufw", "status", "numbered", "verbose").Output()
	if err != nil {
		return nil, ufwInternal("ufw status numbered", err)
	}
	verbose, err := execCommandContext(ctx, "ufw", "status", "verbose").Output()
	if err != nil {
		return nil, ufwInternal("ufw status verbose", err)
	}
	resp := parseUfwStatus(string(numbered) + "\n" + string(verbose))
	return resp, nil
}

var (
	ufwDefaultRe = regexp.MustCompile(`Default:\s*([a-z]+)\s*\(incoming\),\s*([a-z]+)\s*\(outgoing\)`)
	ufwRuleRe    = regexp.MustCompile(`^\[\s*(\d+)\]\s+(.+?)\s{2,}(ALLOW IN|DENY IN|REJECT IN|LIMIT IN|ALLOW OUT|DENY OUT|REJECT OUT)\s+(.+?)\s*$`)
	// portProto matches "22/tcp", "80", "1000:2000/udp"
	ufwPortProtoRe = regexp.MustCompile(`^([\d:]+)(?:/(tcp|udp))?(\s*\(v6\))?$`)
)

func parseUfwStatus(out string) ufwStatusResponse {
	resp := ufwStatusResponse{Rules: []ufwRule{}}
	for _, line := range strings.Split(out, "\n") {
		l := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(l, "Status:"):
			// "Status: active" or "Status: inactive" (GH #2271: a substring
			// test read the second as active too).
			resp.Active = strings.TrimSpace(strings.TrimPrefix(l, "Status:")) == "active"
		case strings.HasPrefix(l, "Default:"):
			if m := ufwDefaultRe.FindStringSubmatch(l); m != nil {
				resp.DefaultIn = m[1]
				resp.DefaultOut = m[2]
			}
		default:
			if m := ufwRuleRe.FindStringSubmatch(l); m != nil {
				num, _ := strconv.Atoi(m[1])
				to := strings.TrimSpace(m[2])
				action := strings.TrimSpace(m[3])
				from := strings.TrimSpace(m[4])
				rule := ufwRule{Num: num, Action: action, From: from, To: to}
				if pm := ufwPortProtoRe.FindStringSubmatch(to); pm != nil {
					rule.Port = pm[1]
					if pm[2] != "" {
						rule.Proto = pm[2]
					}
				}
				resp.Rules = append(resp.Rules, rule)
			}
		}
	}
	return resp
}

// ---- security.ufw.rule.add -------------------------------------------------

type ufwRuleAddParams struct {
	Action string `json:"action"`
	Port   string `json:"port"`
	Proto  string `json:"proto,omitempty"`
	From   string `json:"from,omitempty"`
}

type ufwRuleAddResponse struct {
	Added   bool `json:"added"`
	RuleNum int  `json:"rule_num"`
}

var (
	validUfwActions = map[string]bool{"allow": true, "deny": true, "reject": true}
	validUfwProtos  = map[string]bool{"tcp": true, "udp": true}
	ufwPortRe       = regexp.MustCompile(`^(\d+)(?::(\d+))?$`)
)

func ufwRuleAddHandler(ctx context.Context, params json.RawMessage) (any, error) {
	var p ufwRuleAddParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, ufwInvalidArg(fmt.Sprintf("parse params: %v", err))
	}
	if !validUfwActions[p.Action] {
		return nil, ufwInvalidArg("action must be allow|deny|reject")
	}
	if p.Proto == "" {
		p.Proto = "tcp"
	}
	if !validUfwProtos[p.Proto] {
		return nil, ufwInvalidArg("proto must be tcp|udp")
	}
	pm := ufwPortRe.FindStringSubmatch(p.Port)
	if pm == nil {
		return nil, ufwInvalidArg(`port must be "<N>" or "<N>:<M>"`)
	}
	lo, _ := strconv.Atoi(pm[1])
	if lo < 1 || lo > 65535 {
		return nil, ufwInvalidArg("port out of range 1..65535")
	}
	if pm[2] != "" {
		hi, _ := strconv.Atoi(pm[2])
		if hi < 1 || hi > 65535 || hi <= lo {
			return nil, ufwInvalidArg("port range must be lo<hi within 1..65535")
		}
	}
	if p.From != "" {
		if _, _, err := net.ParseCIDR(p.From); err != nil {
			if ip := net.ParseIP(p.From); ip == nil {
				return nil, ufwInvalidArg("from must be IP or CIDR")
			}
		}
	}
	if isFirewalldBackend() {
		return firewalldRuleAddHandler(ctx, p)
	}
	args := []string{p.Action}
	if p.From != "" {
		args = append(args, "from", p.From, "to", "any", "port", p.Port)
		if p.Proto != "" {
			args = append(args, "proto", p.Proto)
		}
	} else {
		args = append(args, p.Port+"/"+p.Proto)
	}
	if _, err := execCommandContext(ctx, "ufw", args...).CombinedOutput(); err != nil {
		return nil, ufwInternal("ufw "+p.Action+" failed", err)
	}
	num := ufwLookupRuleNum(ctx, p)
	return ufwRuleAddResponse{Added: true, RuleNum: num}, nil
}

func ufwLookupRuleNum(ctx context.Context, p ufwRuleAddParams) int {
	out, err := execCommandContext(ctx, "ufw", "status", "numbered").Output()
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(out), "\n") {
		l := strings.TrimSpace(line)
		if !strings.Contains(l, p.Port) {
			continue
		}
		if m := ufwRuleRe.FindStringSubmatch(l); m != nil {
			n, _ := strconv.Atoi(m[1])
			return n
		}
	}
	return 0
}

// ---- security.ufw.rule.delete ----------------------------------------------

type ufwRuleDeleteParams struct {
	Num int `json:"num"`
}

func ufwRuleDeleteHandler(ctx context.Context, params json.RawMessage) (any, error) {
	var p ufwRuleDeleteParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, ufwInvalidArg(fmt.Sprintf("parse params: %v", err))
	}
	if p.Num < 1 || p.Num > 1000 {
		return nil, ufwInvalidArg("num out of range 1..1000")
	}
	if isFirewalldBackend() {
		return firewalldRuleDeleteHandler(ctx, p.Num)
	}
	cmd := execCommandContext(ctx, "ufw", "--force", "delete", strconv.Itoa(p.Num))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, ufwInternal("ufw delete: "+strings.TrimSpace(stderr.String()), err)
	}
	return map[string]bool{"deleted": true}, nil
}

// ---- security.ufw.default.set ----------------------------------------------

type ufwDefaultSetParams struct {
	Chain  string `json:"chain"`
	Policy string `json:"policy"`
}

var (
	validUfwChains = map[string]bool{"incoming": true, "outgoing": true, "routed": true}
)

func ufwDefaultSetHandler(ctx context.Context, params json.RawMessage) (any, error) {
	var p ufwDefaultSetParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, ufwInvalidArg(fmt.Sprintf("parse params: %v", err))
	}
	if !validUfwChains[p.Chain] {
		return nil, ufwInvalidArg("chain must be incoming|outgoing|routed")
	}
	if !validUfwActions[p.Policy] {
		return nil, ufwInvalidArg("policy must be allow|deny|reject")
	}
	if isFirewalldBackend() {
		return firewalldDefaultSetHandler(ctx, p)
	}
	if _, err := execCommandContext(ctx, "ufw", "default", p.Policy, p.Chain).CombinedOutput(); err != nil {
		return nil, ufwInternal("ufw default failed", err)
	}
	return map[string]bool{"applied": true}, nil
}

// ---- security.ufw.enable / .disable ----------------------------------------

type ufwToggleParams struct {
	Confirm string `json:"confirm"`
}

func ufwEnableHandler(ctx context.Context, params json.RawMessage) (any, error) {
	var p ufwToggleParams
	_ = json.Unmarshal(params, &p)
	if p.Confirm != "YES" {
		return nil, ufwInvalidArg(`confirm must be "YES" — refusing to enable firewall without explicit confirmation`)
	}
	if isFirewalldBackend() {
		if _, err := execCommandContext(ctx, "systemctl", "enable", "--now", "firewalld").CombinedOutput(); err != nil {
			return nil, ufwInternal("firewalld enable failed", err)
		}
		return map[string]bool{"active": true}, nil
	}
	if _, err := execCommandContext(ctx, "ufw", "--force", "enable").CombinedOutput(); err != nil {
		return nil, ufwInternal("ufw enable failed", err)
	}
	return map[string]bool{"active": true}, nil
}

func ufwDisableHandler(ctx context.Context, params json.RawMessage) (any, error) {
	var p ufwToggleParams
	_ = json.Unmarshal(params, &p)
	if p.Confirm != "YES" {
		return nil, ufwInvalidArg(`confirm must be "YES" — refusing to disable firewall without explicit confirmation`)
	}
	if isFirewalldBackend() {
		if _, err := execCommandContext(ctx, "systemctl", "disable", "--now", "firewalld").CombinedOutput(); err != nil {
			return nil, ufwInternal("firewalld disable failed", err)
		}
		return map[string]bool{"active": false}, nil
	}
	if _, err := execCommandContext(ctx, "ufw", "--force", "disable").CombinedOutput(); err != nil {
		return nil, ufwInternal("ufw disable failed", err)
	}
	return map[string]bool{"active": false}, nil
}

// ---- firewalld backend (AlmaLinux / RHEL 8/9/10) ---------------------------

var isFirewalldBackend = func() bool {
	return osinfo.Detect().FirewallBackend() == "firewalld"
}

type firewalldRuleSpec struct {
	rule     ufwRule
	portSpec string // e.g. "80/tcp" or "1000-2000/tcp"
	richRule string // raw rich-rule string for --remove-rich-rule
}

var (
	firewalldRichFromRe   = regexp.MustCompile(`source\s+address="([^"]+)"`)
	firewalldRichPortRe   = regexp.MustCompile(`port\s+port="([\d\-]+)"\s+protocol="(tcp|udp)"`)
	firewalldRichActionRe = regexp.MustCompile(`\b(accept|drop|reject)\s*$`)
)

func firewalldStatusHandler(ctx context.Context) (ufwStatusResponse, error) {
	activeOut, err := execCommandContext(ctx, "systemctl", "is-active", "firewalld").Output()
	if err != nil || strings.TrimSpace(string(activeOut)) != "active" {
		return ufwStatusResponse{
			Active:     false,
			DefaultIn:  "deny",
			DefaultOut: "allow",
			Rules:      []ufwRule{},
		}, nil
	}
	out, err := execCommandContext(ctx, "firewall-cmd", "--list-all").Output()
	if err != nil {
		return ufwStatusResponse{}, ufwInternal("firewall-cmd --list-all", err)
	}
	resp, _ := parseFirewalldListAll(string(out), true)
	return resp, nil
}

func parseFirewalldListAll(out string, active bool) (ufwStatusResponse, []firewalldRuleSpec) {
	resp := ufwStatusResponse{
		Active:     active,
		DefaultIn:  "deny",
		DefaultOut: "allow",
		Rules:      []ufwRule{},
	}
	var specs []firewalldRuleSpec
	inRich := false
	nextNum := 1
	for _, rawLine := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(rawLine)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "target:") {
			inRich = false
			target := strings.ToUpper(strings.TrimSpace(strings.TrimPrefix(trimmed, "target:")))
			switch target {
			case "ACCEPT":
				resp.DefaultIn = "allow"
			case "%%REJECT%%":
				resp.DefaultIn = "reject"
			default:
				resp.DefaultIn = "deny"
			}
			continue
		}
		if strings.HasPrefix(trimmed, "ports:") {
			inRich = false
			for _, token := range strings.Fields(strings.TrimSpace(strings.TrimPrefix(trimmed, "ports:"))) {
				portPart, proto, ok := strings.Cut(token, "/")
				if !ok {
					continue
				}
				displayPort := strings.ReplaceAll(portPart, "-", ":")
				r := ufwRule{
					Num:    nextNum,
					Action: "ALLOW IN",
					From:   "Anywhere",
					To:     displayPort + "/" + proto,
					Proto:  proto,
					Port:   displayPort,
				}
				resp.Rules = append(resp.Rules, r)
				specs = append(specs, firewalldRuleSpec{rule: r, portSpec: token})
				nextNum++
			}
			continue
		}
		if strings.HasPrefix(trimmed, "rich rules:") {
			inRich = true
			continue
		}
		if inRich && strings.HasPrefix(trimmed, "rule ") {
			pm := firewalldRichPortRe.FindStringSubmatch(trimmed)
			if pm == nil {
				continue
			}
			displayPort := strings.ReplaceAll(pm[1], "-", ":")
			proto := pm[2]
			from := "Anywhere"
			if fm := firewalldRichFromRe.FindStringSubmatch(trimmed); fm != nil {
				from = fm[1]
			}
			action := "ALLOW IN"
			if am := firewalldRichActionRe.FindStringSubmatch(trimmed); am != nil {
				switch am[1] {
				case "drop":
					action = "DENY IN"
				case "reject":
					action = "REJECT IN"
				}
			}
			r := ufwRule{
				Num:    nextNum,
				Action: action,
				From:   from,
				To:     displayPort + "/" + proto,
				Proto:  proto,
				Port:   displayPort,
			}
			resp.Rules = append(resp.Rules, r)
			specs = append(specs, firewalldRuleSpec{rule: r, richRule: trimmed})
			nextNum++
			continue
		}
		if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "rule ") {
			inRich = false
		}
	}
	return resp, specs
}

func buildFirewalldRichRule(p ufwRuleAddParams) string {
	family := "ipv4"
	if p.From != "" {
		if ip := net.ParseIP(p.From); ip != nil && ip.To4() == nil {
			family = "ipv6"
		} else if ip, _, err := net.ParseCIDR(p.From); err == nil && ip.To4() == nil {
			family = "ipv6"
		}
	}
	action := "accept"
	switch p.Action {
	case "deny":
		action = "drop"
	case "reject":
		action = "reject"
	}
	fwPort := strings.ReplaceAll(p.Port, ":", "-")
	if p.From != "" {
		return fmt.Sprintf(`rule family="%s" source address="%s" port port="%s" protocol="%s" %s`,
			family, p.From, fwPort, p.Proto, action)
	}
	return fmt.Sprintf(`rule family="%s" port port="%s" protocol="%s" %s`,
		family, fwPort, p.Proto, action)
}

func firewalldRuleAddHandler(ctx context.Context, p ufwRuleAddParams) (any, error) {
	fwPort := strings.ReplaceAll(p.Port, ":", "-")
	if p.Action == "allow" && p.From == "" {
		portArg := fmt.Sprintf("--add-port=%s/%s", fwPort, p.Proto)
		if out, err := execCommandContext(ctx, "firewall-cmd", "--permanent", portArg).CombinedOutput(); err != nil {
			return nil, ufwInternal("firewall-cmd "+portArg+" ("+strings.TrimSpace(string(out))+")", err)
		}
	} else {
		rich := buildFirewalldRichRule(p)
		if out, err := execCommandContext(ctx, "firewall-cmd", "--permanent", "--add-rich-rule="+rich).CombinedOutput(); err != nil {
			return nil, ufwInternal("firewall-cmd add rich rule ("+strings.TrimSpace(string(out))+")", err)
		}
	}
	if out, err := execCommandContext(ctx, "firewall-cmd", "--reload").CombinedOutput(); err != nil {
		return nil, ufwInternal("firewall-cmd --reload ("+strings.TrimSpace(string(out))+")", err)
	}
	return ufwRuleAddResponse{Added: true, RuleNum: 0}, nil
}

func firewalldRuleDeleteHandler(ctx context.Context, num int) (any, error) {
	out, err := execCommandContext(ctx, "firewall-cmd", "--list-all").Output()
	if err != nil {
		return nil, ufwInternal("firewall-cmd --list-all", err)
	}
	_, specs := parseFirewalldListAll(string(out), true)
	if num < 1 || num > len(specs) {
		return nil, ufwInvalidArg(fmt.Sprintf("rule #%d not found", num))
	}
	target := specs[num-1]
	var arg string
	if target.portSpec != "" {
		arg = "--remove-port=" + target.portSpec
	} else {
		arg = "--remove-rich-rule=" + target.richRule
	}
	if cout, err := execCommandContext(ctx, "firewall-cmd", "--permanent", arg).CombinedOutput(); err != nil {
		return nil, ufwInternal("firewall-cmd delete ("+strings.TrimSpace(string(cout))+")", err)
	}
	if cout, err := execCommandContext(ctx, "firewall-cmd", "--reload").CombinedOutput(); err != nil {
		return nil, ufwInternal("firewall-cmd --reload ("+strings.TrimSpace(string(cout))+")", err)
	}
	return map[string]bool{"deleted": true}, nil
}

func firewalldDefaultSetHandler(ctx context.Context, p ufwDefaultSetParams) (any, error) {
	if p.Chain != "incoming" {
		return map[string]bool{"applied": true}, nil
	}
	target := "DROP"
	switch p.Policy {
	case "allow":
		target = "ACCEPT"
	case "reject":
		target = "%%REJECT%%"
	}
	if cout, err := execCommandContext(ctx, "firewall-cmd", "--permanent", "--set-target="+target).CombinedOutput(); err != nil {
		return nil, ufwInternal("firewall-cmd --set-target ("+strings.TrimSpace(string(cout))+")", err)
	}
	if cout, err := execCommandContext(ctx, "firewall-cmd", "--reload").CombinedOutput(); err != nil {
		return nil, ufwInternal("firewall-cmd --reload ("+strings.TrimSpace(string(cout))+")", err)
	}
	return map[string]bool{"applied": true}, nil
}

func init() {
	Default.Register("security.ufw.status", ufwStatusHandler)
	Default.Register("security.ufw.rule.add", ufwRuleAddHandler)
	Default.Register("security.ufw.rule.delete", ufwRuleDeleteHandler)
	Default.Register("security.ufw.default.set", ufwDefaultSetHandler)
	Default.Register("security.ufw.enable", ufwEnableHandler)
	Default.Register("security.ufw.disable", ufwDisableHandler)
}
