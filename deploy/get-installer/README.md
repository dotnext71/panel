# get.dotnext-panel.com — short installer URL (GH #444)

> **Status: not deployed.** The `dotnext-panel.com` domain is not registered
> yet, so nothing is served here today. Until it is, the supported install
> command is the raw GitHub URL (`https://raw.githubusercontent.com/dotnext71/panel/main/bootstrap.sh`,
> or `.../main/install.sh` for the non-TUI path) — every doc and the README use
> that form. This directory is the ready-made setup for when the domain exists.

Gives operators a short, memorable install command once the domain is live:

```bash
curl -fsSL https://get.dotnext-panel.com | bash
```

It resolves to the canonical installer (`install.sh` on `main`). `curl -fsSL`
already follows redirects (`-L`), so a 302 works; a 200-rewrite serves the
script body directly with no hop.

## Set it up (one-time, operator/infra — not shipped by the panel itself)

1. **DNS:** point `get.dotnext-panel.com` at the host/CDN that will serve the
   redirect (A/AAAA, or CNAME to a Pages project).
2. **Serve it** — pick one:
   - **Cloudflare Pages / Netlify:** deploy this directory; `_redirects` rewrites
     every path to the raw `install.sh` (200, body served directly). TLS is
     automatic.
   - **nginx:** install `nginx-get.dotnext-panel.com.conf`, issue a cert
     (`certbot --nginx -d get.dotnext-panel.com`), reload. It 302s to the raw URL.
3. **Verify:** `curl -fsSL https://get.dotnext-panel.com | head` prints the
   `#!/usr/bin/env bash` shebang + the installer header.

## Notes

- The canonical raw URL stays valid as a fallback:
  `https://raw.githubusercontent.com/dotnext71/panel/main/install.sh`
- Channel-aware later (GH #445): point `get.` at a small bootstrapper that picks
  the stable vs development installer once channels exist.
