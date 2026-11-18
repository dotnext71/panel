// DotNextTitle — brand lockup (SVG mark + wordmark) for Header and mobile Drawer.
import { Typography } from "antd";
import { useThemeMode } from "../theme/ThemeModeContext";
import { logoURL, useBranding } from "../hooks/useBranding";

interface DotNextTitleProps {
  collapsed?: boolean;
  text?: string;
  showWordmark?: boolean;
}

export function DotNextTitle({
  collapsed = false,
  text,
  showWordmark = true,
}: DotNextTitleProps) {
  const { mode } = useThemeMode();
  const { brandText, hasLogoLight, hasLogoDark } = useBranding();
  const variant: "light" | "dark" = mode === "dark" ? "dark" : "light";
  const hasCustom = variant === "dark" ? hasLogoDark : hasLogoLight;
  const src = logoURL(variant, hasCustom);
  const customLabel = text ?? brandText;

  return (
    <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
      <img
        src={src}
        alt={customLabel || "dotNext"}
        style={{ height: 30, width: "auto", flexShrink: 0 }}
      />
      {!collapsed && showWordmark && customLabel ? (
        <Typography.Title level={4} style={{ margin: 0, letterSpacing: "-0.02em" }}>
          {customLabel}
        </Typography.Title>
      ) : null}
    </div>
  );
}

