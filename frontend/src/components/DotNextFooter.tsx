// DotNextFooter — compact status/attribution footer rendered at the bottom of both shells.
import { useTranslation } from "react-i18next";
import { GithubOutlined } from "@icons";
import { Grid, Layout, Space, Tag, Typography } from "antd";

import { useThemeMode } from "../theme/ThemeModeContext";
import pkg from "../../package.json";

const SOURCE_URL = "https://github.com/dotnext71/panel";

export function DotNextFooter() {
  const { t } = useTranslation();
  const { mode } = useThemeMode();
  const screens = Grid.useBreakpoint();
  const isWide = screens.lg !== false;
  const showExtras = screens.sm !== false;
  const logoSrc =
    mode === "dark" ? "/images/dotnext_logo_dark.svg" : "/images/dotnext_logo.svg";

  return (
    <Layout.Footer
      style={{
        display: "flex",
        flexDirection: isWide ? "row" : "column",
        alignItems: "center",
        justifyContent: isWide ? "space-between" : "center",
        textAlign: isWide ? "left" : "center",
        gap: isWide ? 16 : 8,
        padding: isWide ? "12px 24px" : "10px 12px",
        background: "transparent",
        borderTop: "1px solid var(--dn-border)",
      }}
    >
      <Space size={12}>
        <img
          src={logoSrc}
          alt="dotNext"
          style={{ height: 24, width: "auto", flexShrink: 0 }}
        />
        {showExtras && (
          <Typography.Text type="secondary" style={{ fontSize: 12 }}>
            Multi-Tenant Hosting Control Plane
          </Typography.Text>
        )}
      </Space>

      <Space size={10}>
        <a
          href={SOURCE_URL}
          target="_blank"
          rel="noreferrer"
          aria-label={t("dotnextfooter.source_code", "Source code")}
          style={{ color: "inherit", display: "inline-flex", alignItems: "center" }}
        >
          <GithubOutlined style={{ fontSize: 16 }} />
        </a>
        {showExtras && (
          <>
            <Typography.Text type="secondary">·</Typography.Text>
            <Typography.Text type="secondary" style={{ fontSize: 12 }}>
              <a
                href="https://opensource.org/license/mit"
                target="_blank"
                rel="noreferrer"
                style={{ color: "inherit" }}
              >
                MIT
              </a>
            </Typography.Text>
          </>
        )}
        <Tag
          className="dn-tabular"
          style={{ marginInlineEnd: 0, fontSize: 11, lineHeight: "18px" }}
        >
          v{pkg.version}
        </Tag>
      </Space>
    </Layout.Footer>
  );
}
