package paths

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestPathsAreAbsoluteAndRebranded(t *testing.T) {
	t.Parallel()

	absPaths := []string{
		PanelBinary,
		APIBinary,
		CLIBinary,
		AgentBinary,
		WebDAVBinary,
		SSHShellBinary,
		MailhookBinary,
		SendmailBinary,
		InstallerBinary,
		ConfigDir,
		ConfigFile,
		EnvFile,
		SSOKeyFile,
		PDNSEnvFile,
		FPMConfigRoot,
		PHPVerPinRoot,
		PanelTLSCert,
		PanelTLSKey,
		PanelRunDir,
		PanelSocket,
		SSOSocket,
		AgentRunDir,
		AgentSocket,
		AgentPTYSocket,
		KratosPublicSock,
		KratosAdminSock,
		WebDAVRunDir,
		RepoDir,
		WebmailDir,
		DataDir,
		BackupsDir,
		AgentStagingDir,
		FTPJailsDir,
		LogDir,
		TerminalLogDir,
		ShareDir,
	}

	for _, p := range absPaths {
		if !filepath.IsAbs(p) {
			t.Errorf("expected absolute path, got %q", p)
		}
		if strings.Contains(strings.ToLower(p), "jabali") {
			t.Errorf("canonical dotNext path %q must not contain legacy branding", p)
		}
	}

	if AdminBasePath != "/admin" || UserBasePath != "/panel" {
		t.Errorf("unexpected shell paths: %q, %q", AdminBasePath, UserBasePath)
	}
	if APITokenPrefix != "dnt_" || AutomationAuthHdr != "DotNext-HMAC" {
		t.Errorf("unexpected auth identifiers: %q, %q", APITokenPrefix, AutomationAuthHdr)
	}
}
