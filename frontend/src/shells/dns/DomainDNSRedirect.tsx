// GH #1920 (johnnyq): DNS records are managed only under DNS > Zones. The old
// per-domain URLs — /admin/domains/:id/dns and the tenant Web Domain
// page's DNS tab at /panel/domains/:id/dns — stay valid for bookmarks
// and links, and land on the zone's records page in the same shell.
import { Navigate, useLocation, useParams } from "react-router";

export function DomainDNSRedirect() {
  const { id = "" } = useParams<{ id: string }>();
  const { pathname } = useLocation();
  const shell = pathname.startsWith("/admin") ? "/admin" : "/panel";
  return <Navigate to={`${shell}/dns/${encodeURIComponent(id)}`} replace />;
}
