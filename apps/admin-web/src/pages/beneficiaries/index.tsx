import { List, useTable, ShowButton, Show } from "@refinedev/antd";
import { Table, Space, Tag, Card, Descriptions } from "antd";
import { useShow } from "@refinedev/core";

export const BeneficiaryList: React.FC = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  const statusColor = (status: string) => {
    const colors: Record<string, string> = {
      active: "green",
      inactive: "default",
      pending_verification: "orange",
      blocked: "red",
    };
    return colors[status] || "default";
  };

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column
          dataIndex="first_name"
          title="Name"
          render={(_, record: any) =>
            `${record.first_name} ${record.last_name}`
          }
        />
        <Table.Column dataIndex="relationship" title="Relationship" />
        <Table.Column dataIndex="phone_primary" title="Phone" />
        <Table.Column dataIndex="city" title="City" />
        <Table.Column dataIndex="country" title="Country" />
        <Table.Column
          dataIndex="preferred_delivery_method"
          title="Delivery"
          render={(method) => method?.replace("_", " ")}
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
          dataIndex="transfer_count"
          title="Transfers"
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

export const BeneficiaryShow: React.FC = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Card title="Beneficiary Details" style={{ marginBottom: 16 }}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Full Name">
            {record?.first_name} {record?.middle_name} {record?.last_name}
          </Descriptions.Item>
          <Descriptions.Item label="Nickname">{record?.nickname || "-"}</Descriptions.Item>
          <Descriptions.Item label="Relationship">{record?.relationship}</Descriptions.Item>
          <Descriptions.Item label="Status">{record?.status}</Descriptions.Item>
          <Descriptions.Item label="Primary Phone">{record?.phone_primary}</Descriptions.Item>
          <Descriptions.Item label="Secondary Phone">{record?.phone_secondary || "-"}</Descriptions.Item>
          <Descriptions.Item label="Email">{record?.email || "-"}</Descriptions.Item>
          <Descriptions.Item label="Favorite">{record?.is_favorite ? "Yes" : "No"}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="Address" style={{ marginBottom: 16 }}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Address">
            {record?.address_line1} {record?.address_line2}
          </Descriptions.Item>
          <Descriptions.Item label="City">{record?.city}</Descriptions.Item>
          <Descriptions.Item label="State">{record?.state}</Descriptions.Item>
          <Descriptions.Item label="Country">{record?.country}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="Payment Details" style={{ marginBottom: 16 }}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Bank Name">{record?.bank_name || "-"}</Descriptions.Item>
          <Descriptions.Item label="Account Number">{record?.account_number || "-"}</Descriptions.Item>
          <Descriptions.Item label="Account Name">{record?.account_name || "-"}</Descriptions.Item>
          <Descriptions.Item label="Mobile Wallet">{record?.mobile_wallet_provider || "-"}</Descriptions.Item>
          <Descriptions.Item label="Wallet Number">{record?.mobile_wallet_number || "-"}</Descriptions.Item>
          <Descriptions.Item label="Preferred Method">{record?.preferred_delivery_method}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="Statistics">
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Total Transfers">{record?.transfer_count}</Descriptions.Item>
          <Descriptions.Item label="Total Transferred">
            {record?.preferred_currency} {parseFloat(record?.total_transferred || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Last Transfer">
            {record?.last_transfer_at && new Date(record.last_transfer_at).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Verified">
            {record?.verification_status === "verified" ? "Yes" : "No"}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </Show>
  );
};
