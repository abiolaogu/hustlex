import { List, useTable, ShowButton, Show } from "@refinedev/antd";
import { Table, Space, Tag, Card, Descriptions, Steps } from "antd";
import { useShow } from "@refinedev/core";

export const RemittanceList: React.FC = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  const statusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: "orange",
      quoted: "blue",
      initiated: "cyan",
      processing: "purple",
      in_transit: "geekblue",
      delivered: "lime",
      completed: "green",
      failed: "red",
      cancelled: "default",
      refunded: "magenta",
      on_hold: "volcano",
    };
    return colors[status] || "default";
  };

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="reference" title="Reference" />
        <Table.Column
          dataIndex={["user", "profile", "first_name"]}
          title="Sender"
          render={(_, record: any) =>
            `${record.user?.profile?.first_name || ""} ${record.user?.profile?.last_name || ""}`
          }
        />
        <Table.Column
          dataIndex={["beneficiary", "first_name"]}
          title="Beneficiary"
          render={(_, record: any) =>
            `${record.beneficiary?.first_name || ""} ${record.beneficiary?.last_name || ""}`
          }
        />
        <Table.Column
          dataIndex="source_amount"
          title="Sent"
          render={(amount, record: any) =>
            `${record.source_currency} ${parseFloat(amount || 0).toLocaleString()}`
          }
        />
        <Table.Column
          dataIndex="target_amount"
          title="Received"
          render={(amount, record: any) =>
            `${record.target_currency} ${parseFloat(amount || 0).toLocaleString()}`
          }
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
          dataIndex="created_at"
          title="Date"
          render={(date) => new Date(date).toLocaleDateString()}
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

export const RemittanceShow: React.FC = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;
  const record = data?.data;

  const getStepStatus = (status: string) => {
    const steps = ["pending", "quoted", "initiated", "processing", "in_transit", "delivered", "completed"];
    const current = steps.indexOf(status);
    return current;
  };

  return (
    <Show isLoading={isLoading}>
      <Card title="Remittance Status" style={{ marginBottom: 16 }}>
        <Steps
          current={getStepStatus(record?.status)}
          items={[
            { title: "Pending" },
            { title: "Quoted" },
            { title: "Initiated" },
            { title: "Processing" },
            { title: "In Transit" },
            { title: "Delivered" },
            { title: "Completed" },
          ]}
        />
      </Card>

      <Card title="Remittance Details" style={{ marginBottom: 16 }}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Reference">{record?.reference}</Descriptions.Item>
          <Descriptions.Item label="External Ref">{record?.external_reference || "-"}</Descriptions.Item>
          <Descriptions.Item label="Sender">
            {record?.user?.profile?.first_name} {record?.user?.profile?.last_name}
          </Descriptions.Item>
          <Descriptions.Item label="Beneficiary">
            {record?.beneficiary?.first_name} {record?.beneficiary?.last_name}
          </Descriptions.Item>
          <Descriptions.Item label="Source Amount">
            {record?.source_currency} {parseFloat(record?.source_amount || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Target Amount">
            {record?.target_currency} {parseFloat(record?.target_amount || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="FX Rate">{record?.fx_rate}</Descriptions.Item>
          <Descriptions.Item label="Total Fee">
            {record?.source_currency} {parseFloat(record?.total_fee || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Purpose">{record?.purpose?.replace("_", " ")}</Descriptions.Item>
          <Descriptions.Item label="Delivery Method">{record?.delivery_method?.replace("_", " ")}</Descriptions.Item>
          <Descriptions.Item label="Created">{new Date(record?.created_at).toLocaleString()}</Descriptions.Item>
          <Descriptions.Item label="Estimated Delivery">
            {record?.estimated_delivery && new Date(record.estimated_delivery).toLocaleString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="Compliance">
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Compliance Status">{record?.compliance_status}</Descriptions.Item>
          <Descriptions.Item label="AML Checked">{record?.aml_checked ? "Yes" : "No"}</Descriptions.Item>
          <Descriptions.Item label="AML Checked At">
            {record?.aml_checked_at && new Date(record.aml_checked_at).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Notes">{record?.compliance_notes || "-"}</Descriptions.Item>
        </Descriptions>
      </Card>
    </Show>
  );
};
