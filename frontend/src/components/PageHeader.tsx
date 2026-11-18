import type { ReactNode } from "react";
import { Space, Typography } from "antd";

interface PageHeaderProps {
  icon?: ReactNode;
  title: ReactNode;
  subtitle?: ReactNode;
  actions?: ReactNode;
}

export function PageHeader({ icon, title, subtitle, actions }: PageHeaderProps) {
  return (
    <div
      style={{
        display: "flex",
        flexWrap: "wrap",
        alignItems: "center",
        justifyContent: "space-between",
        gap: 12,
        paddingBottom: 14,
        marginBottom: 16,
        borderBottom: "1px solid var(--dn-border)",
      }}
    >
      <div style={{ minWidth: 0 }}>
        <Typography.Title level={3} style={{ margin: 0, display: "flex", alignItems: "center", gap: 10 }}>
          {icon}
          <span>{title}</span>
        </Typography.Title>
        {subtitle && (
          <Typography.Text type="secondary" style={{ display: "block", marginTop: 4, fontSize: 13 }}>
            {subtitle}
          </Typography.Text>
        )}
      </div>
      {actions && <Space wrap>{actions}</Space>}
    </div>
  );
}
