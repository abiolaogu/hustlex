import { List, useTable, ShowButton, Show } from "@refinedev/antd";
import { Table, Space, Tag, Card, Descriptions, Timeline } from "antd";
import { useShow } from "@refinedev/core";

export const BookingList: React.FC = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  const statusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: "orange",
      confirmed: "blue",
      paid: "cyan",
      in_progress: "purple",
      completed: "green",
      cancelled: "red",
      refunded: "magenta",
    };
    return colors[status] || "default";
  };

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="reference" title="Reference" />
        <Table.Column
          dataIndex={["service", "title"]}
          title="Service"
        />
        <Table.Column
          dataIndex={["consumer", "profile", "first_name"]}
          title="Consumer"
          render={(_, record: any) =>
            `${record.consumer?.profile?.first_name || ""} ${record.consumer?.profile?.last_name || ""}`
          }
        />
        <Table.Column
          dataIndex={["provider", "profile", "first_name"]}
          title="Provider"
          render={(_, record: any) =>
            `${record.provider?.profile?.first_name || ""} ${record.provider?.profile?.last_name || ""}`
          }
        />
        <Table.Column
          dataIndex="total_amount"
          title="Amount"
          render={(amount, record: any) =>
            `${record.currency || "NGN"} ${parseFloat(amount || 0).toLocaleString()}`
          }
        />
        <Table.Column
          dataIndex="scheduled_date"
          title="Date"
          render={(date) => date && new Date(date).toLocaleDateString()}
        />
        <Table.Column
          dataIndex="status"
          title="Status"
          render={(status) => (
            <Tag color={statusColor(status)}>
              {status?.replace("_", " ").toUpperCase()}
            </Tag>
          )}
        />
        <Table.Column
          title="Actions"
          dataIndex="actions"
          render={(_, record: any) => (
            <Space>
              <ShowButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};

export const BookingShow: React.FC = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Card title="Booking Details" style={{ marginBottom: 16 }}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Reference">{record?.reference}</Descriptions.Item>
          <Descriptions.Item label="Status">{record?.status}</Descriptions.Item>
          <Descriptions.Item label="Service">{record?.service?.title}</Descriptions.Item>
          <Descriptions.Item label="Amount">
            {record?.currency} {parseFloat(record?.total_amount || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Consumer">
            {record?.consumer?.profile?.first_name} {record?.consumer?.profile?.last_name}
          </Descriptions.Item>
          <Descriptions.Item label="Provider">
            {record?.provider?.profile?.first_name} {record?.provider?.profile?.last_name}
          </Descriptions.Item>
          <Descriptions.Item label="Scheduled Date">
            {record?.scheduled_date && new Date(record.scheduled_date).toLocaleDateString()}
          </Descriptions.Item>
          <Descriptions.Item label="Time">
            {record?.scheduled_time_start} - {record?.scheduled_time_end}
          </Descriptions.Item>
          <Descriptions.Item label="Address" span={2}>
            {record?.service_address}, {record?.service_city}, {record?.service_state}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="Timeline">
        <Timeline
          items={[
            { children: `Created: ${new Date(record?.created_at).toLocaleString()}` },
            record?.paid_at && { children: `Paid: ${new Date(record.paid_at).toLocaleString()}` },
            record?.started_at && { children: `Started: ${new Date(record.started_at).toLocaleString()}` },
            record?.completed_at && { children: `Completed: ${new Date(record.completed_at).toLocaleString()}`, color: "green" },
            record?.cancelled_at && { children: `Cancelled: ${new Date(record.cancelled_at).toLocaleString()}`, color: "red" },
          ].filter(Boolean) as any}
        />
      </Card>
    </Show>
  );
};
