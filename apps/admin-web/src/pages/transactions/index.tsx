import { List, useTable, ShowButton, Show } from "@refinedev/antd";
import { Table, Space, Tag, Card, Descriptions } from "antd";
import { useShow } from "@refinedev/core";

export const TransactionList: React.FC = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  const statusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: "orange",
      processing: "blue",
      completed: "green",
      failed: "red",
      reversed: "magenta",
    };
    return colors[status] || "default";
  };

  const typeColor = (type: string) => {
    const colors: Record<string, string> = {
      credit: "green",
      debit: "red",
      transfer: "blue",
      payment: "purple",
      refund: "cyan",
      withdrawal: "orange",
      deposit: "lime",
    };
    return colors[type] || "default";
  };

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="reference" title="Reference" />
        <Table.Column
          dataIndex="type"
          title="Type"
          render={(type) => (
            <Tag color={typeColor(type)}>{type?.toUpperCase()}</Tag>
          )}
        />
        <Table.Column
          dataIndex="amount"
          title="Amount"
          render={(amount, record: any) =>
            `${record.currency || "NGN"} ${parseFloat(amount || 0).toLocaleString()}`
          }
        />
        <Table.Column
          dataIndex="fee"
          title="Fee"
          render={(fee, record: any) =>
            fee > 0 ? `${record.currency || "NGN"} ${parseFloat(fee).toLocaleString()}` : "-"
          }
        />
        <Table.Column
          dataIndex="status"
          title="Status"
          render={(status) => (
            <Tag color={statusColor(status)}>{status?.toUpperCase()}</Tag>
          )}
        />
        <Table.Column
          dataIndex="created_at"
          title="Date"
          render={(date) => new Date(date).toLocaleString()}
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

export const TransactionShow: React.FC = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Card title="Transaction Details">
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Reference">{record?.reference}</Descriptions.Item>
          <Descriptions.Item label="External Reference">{record?.external_reference || "-"}</Descriptions.Item>
          <Descriptions.Item label="Type">{record?.type}</Descriptions.Item>
          <Descriptions.Item label="Status">{record?.status}</Descriptions.Item>
          <Descriptions.Item label="Amount">
            {record?.currency} {parseFloat(record?.amount || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Fee">
            {record?.currency} {parseFloat(record?.fee || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Balance Before">
            {record?.currency} {parseFloat(record?.balance_before || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Balance After">
            {record?.currency} {parseFloat(record?.balance_after || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Description" span={2}>
            {record?.description}
          </Descriptions.Item>
          <Descriptions.Item label="Initiated">
            {record?.initiated_at && new Date(record.initiated_at).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Completed">
            {record?.completed_at && new Date(record.completed_at).toLocaleString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </Show>
  );
};
