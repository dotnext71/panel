import { useQuery } from "@tanstack/react-query";
import { Button, Card, Empty, Table, Tag, Typography } from "antd";
import { ReloadOutlined } from "@icons";

import { apiClient } from "../../../apiClient";
import { Timestamp } from "../../../components/Timestamp";

export interface SupervisedWorker {
  name: string;
  running: boolean;
  restarts: number;
  last_start?: string;
  last_error?: string;
}

interface WorkersResponse {
  data: SupervisedWorker[];
  total: number;
}

export function SupervisedWorkersCard() {
  const { data, isLoading, isError, refetch, isFetching } = useQuery<WorkersResponse>({
    queryKey: ["admin", "system", "workers"],
    queryFn: async () => (await apiClient.get<WorkersResponse>("/admin/system/workers")).data,
    refetchInterval: 15_000,
    retry: 1,
  });

  const rows = data?.data ?? [];
  const runningCount = rows.filter((w) => w.running).length;

  return (
    <Card
      size="small"
      title={
        <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
          <span>Supervised Workers</span>
          {rows.length > 0 && (
            <Tag color={runningCount === rows.length ? "green" : "orange"} className="dn-tabular">
              {runningCount}/{rows.length} running
            </Tag>
          )}
        </div>
      }
      extra={
        <Button
          size="small"
          icon={<ReloadOutlined />}
          onClick={() => void refetch()}
          loading={isFetching && !isLoading}
        >
          Refresh
        </Button>
      }
    >
      {isError ? (
        <Typography.Text type="secondary">
          Worker telemetry endpoint is unavailable on this host.
        </Typography.Text>
      ) : (
        <Table<SupervisedWorker>
          size="small"
          rowKey="name"
          loading={isLoading}
          dataSource={rows}
          pagination={false}
          scroll={{ x: "max-content" }}
          locale={{ emptyText: <Empty description="No background workers registered" /> }}
          columns={[
            {
              title: "Worker",
              dataIndex: "name",
              render: (name: string) => <Typography.Text code>{name}</Typography.Text>,
            },
            {
              title: "State",
              dataIndex: "running",
              width: 100,
              render: (running: boolean) =>
                running ? <Tag color="green">running</Tag> : <Tag color="default">stopped</Tag>,
            },
            {
              title: "Restarts",
              dataIndex: "restarts",
              width: 90,
              className: "dn-tabular",
            },
            {
              title: "Last Start",
              dataIndex: "last_start",
              render: (ts?: string) => (ts ? <Timestamp value={ts} /> : "—"),
            },
            {
              title: "Last Error",
              dataIndex: "last_error",
              ellipsis: true,
              render: (err?: string) =>
                err ? <Typography.Text type="danger">{err}</Typography.Text> : "—",
            },
          ]}
        />
      )}
    </Card>
  );
}
