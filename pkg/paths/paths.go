// SPDX-License-Identifier: MIT

// Package paths defines the canonical filesystem, socket, systemd unit,
// database, and protocol-header identifiers for dotNext.
//
// Centralizing these values prevents drift across dotnext-api, dotnext-agent,
// dotnext-cli, dotnext-installer, and runtime configuration loaders.
package paths

const (
	// Product identity
	ProductName = "dotNext"
	ProductSlug = "dotnext"

	// Binaries (installed in /usr/local/bin)
	BinDir          = "/usr/local/bin"
	PanelBinary     = "/usr/local/bin/dotnext-api"
	APIBinary       = "/usr/local/bin/dotnext-api"
	CLIBinary       = "/usr/local/bin/dotnext-cli"
	AgentBinary     = "/usr/local/bin/dotnext-agent"
	WebDAVBinary    = "/usr/local/bin/dotnext-webdav"
	SSHShellBinary  = "/usr/local/bin/dotnext-ssh-shell"
	MailhookBinary  = "/usr/local/bin/dotnext-mailhook"
	SendmailBinary  = "/usr/local/bin/dotnext-sendmail"
	InstallerBinary = "/usr/local/bin/dotnext-installer"

	// Service users & groups
	ServiceUser  = "dotnext"
	ServiceGroup = "dotnext"
	SocketGroup  = "dotnext-sockets"

	// Systemd service units
	PanelService     = "dotnext-api"
	AgentService     = "dotnext-agent"
	KratosService    = "dotnext-kratos"
	MailRelayService = "dotnext-mailrelay"
	WebDAVTemplate   = "dotnext-webdav@"

	// Configuration & state directories
	ConfigDir         = "/etc/dotnext"
	ConfigFile        = "/etc/dotnext/config.toml"
	EnvFile           = "/etc/dotnext/panel.env"
	SSOKeyFile        = "/etc/dotnext/sso.key"
	PDNSEnvFile       = "/etc/dotnext/pdns.env"
	FPMConfigRoot     = "/etc/dotnext/fpm"
	PHPVerPinRoot     = "/etc/dotnext/user-phpver"
	BulwarkJWTSecret  = "/etc/dotnext/bulwark-jwt-auth.secret"
	DeployProfileFile = "/etc/dotnext/deploy-profile"
	TLSKeyDir         = "/etc/dotnext/tls"
	PanelTLSCert      = "/etc/dotnext/tls/panel.crt"
	PanelTLSKey       = "/etc/dotnext/tls/panel.key"

	// Runtime Unix sockets
	PanelRunDir       = "/run/dotnext-panel"
	PanelSocket       = "/run/dotnext-panel/api.sock"
	SSOSocket         = "/run/dotnext-panel/sso.sock"
	AgentRunDir       = "/run/dotnext"
	AgentSocket       = "/run/dotnext/agent.sock"
	AgentPTYSocket    = "/run/dotnext/agent-pty.sock"
	KratosRunDir      = "/run/dotnext-kratos"
	KratosPublicSock  = "/run/dotnext-kratos/public.sock"
	KratosAdminSock   = "/run/dotnext-kratos/admin.sock"
	WebDAVRunDir      = "/run/dotnext-webdav"

	// Persistent data & staging directories
	RepoDir         = "/opt/dotnext"
	WebmailDir      = "/opt/dotnext-webmail"
	DataDir         = "/var/lib/dotnext"
	BackupsDir      = "/var/lib/dotnext-backups"
	AgentStagingDir = "/var/lib/dotnext-agent-staging"
	FTPJailsDir     = "/var/lib/dotnext-ftp-jails"
	LogDir          = "/var/log/dotnext"
	TerminalLogDir  = "/var/log/dotnext/terminal"
	ShareDir        = "/usr/local/share/dotnext"

	// Databases
	PanelDBName = "dotnext_panel"
	PDNSDBName  = "dotnext_pdns"

	// Auth, wire headers, and Redis keys
	APITokenPrefix    = "dnt_"
	AutomationAuthHdr = "DotNext-HMAC"
	ImpersonateHeader = "X-DotNext-Act-As"
	WebhookSigHeader  = "X-DotNext-Signature"
	AuditRedisStream  = "dotnext:audit:queue"
	NotifyRedisPrefix = "dotnext:notifications:"

	// Primary SPA shell URL prefixes
	AdminBasePath = "/admin"
	UserBasePath  = "/panel"
)
