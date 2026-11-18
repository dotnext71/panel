// GH #1920 (johnnyq): DNS records live under DNS > Zones. The old per-domain
// URLs (/admin/domains/:id/dns, and the tenant Web Domain page's DNS tab
// at /panel/domains/:id/dns) stay valid for bookmarks and links, and
// land on the zone's records page in the same shell.
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes, useLocation } from "react-router";
import { describe, expect, it } from "vitest";

import { DomainDNSRedirect } from "./DomainDNSRedirect";

function Here() {
  const { pathname } = useLocation();
  return <div>at:{pathname}</div>;
}

const renderAt = (url: string) =>
  render(
    <MemoryRouter initialEntries={[url]}>
      <Routes>
        <Route path="/admin/domains/:id/dns" element={<DomainDNSRedirect />} />
        <Route path="/panel/domains/:id/dns" element={<DomainDNSRedirect />} />
        <Route path="/admin/dns/:id" element={<Here />} />
        <Route path="/panel/dns/:id" element={<Here />} />
      </Routes>
    </MemoryRouter>,
  );

describe("DomainDNSRedirect (GH #1920)", () => {
  it.each([
    ["/admin/domains/d1/dns", "/admin/dns/d1"],
    ["/panel/domains/d1/dns", "/panel/dns/d1"],
  ])("sends %s to %s", (from, to) => {
    renderAt(from);
    expect(screen.getByText(`at:${to}`)).toBeInTheDocument();
  });
});
