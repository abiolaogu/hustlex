import { List, useTable, ShowButton, Show } from "@refinedev/antd";
import { Table, Space, Tag, Card, Descriptions, Progress, List as AntList } from "antd";
import { useShow } from "@refinedev/core";

export const SavingsCircleList: React.FC = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  const statusColor = (status: string) => {
    const colors: Record<string, string> = {
      forming: "blue",
      active: "green",
      paused: "orange",
      completed: "cyan",
      dissolved: "red",
    };
    return colors[status] || "default";
  };

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="name" title="Name" />
        <Table.Column
          dataIndex={["creator", "profile", "first_name"]}
          title="Creator"
          render={(_, record: any) =>
            `${record.creator?.profile?.first_name || ""} ${record.creator?.profile?.last_name || ""}`
          }
        />
        <Table.Column
          dataIndex="contribution_amount"
          title="Contribution"
          render={(amount, record: any) =>
            `${record.currency || "NGN"} ${parseFloat(amount || 0).toLocaleString()}`
          }
        />
        <Table.Column dataIndex="contribution_frequency" title="Frequency" />
        <Table.Column
          dataIndex="max_members"
          title="Members"
          render={(max, record: any) =>
            `${record.members?.length || 0}/${max}`
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

export const SavingsCircleShow: React.FC = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;
  const record = data?.data;

  const memberProgress = ((record?.current_cycle || 0) / (record?.total_cycles || 1)) * 100;

  return (
    <Show isLoading={isLoading}>
      <Card title="Circle Details" style={{ marginBottom: 16 }}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Name">{record?.name}</Descriptions.Item>
          <Descriptions.Item label="Status">{record?.status}</Descriptions.Item>
          <Descriptions.Item label="Creator">
            {record?.creator?.profile?.first_name} {record?.creator?.profile?.last_name}
          </Descriptions.Item>
          <Descriptions.Item label="Description">{record?.description}</Descriptions.Item>
          <Descriptions.Item label="Contribution">
            {record?.currency} {parseFloat(record?.contribution_amount || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Frequency">{record?.contribution_frequency}</Descriptions.Item>
          <Descriptions.Item label="Start Date">
            {record?.start_date && new Date(record.start_date).toLocaleDateString()}
          </Descriptions.Item>
          <Descriptions.Item label="Next Contribution">
            {record?.next_contribution_date && new Date(record.next_contribution_date).toLocaleDateString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="Progress" style={{ marginBottom: 16 }}>
        <Progress percent={memberProgress} status="active" />
        <Descriptions bordered column={2} style={{ marginTop: 16 }}>
          <Descriptions.Item label="Current Cycle">{record?.current_cycle}</Descriptions.Item>
          <Descriptions.Item label="Total Cycles">{record?.total_cycles}</Descriptions.Item>
          <Descriptions.Item label="Total Contributed">
            {record?.currency} {parseFloat(record?.total_contributed || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Total Disbursed">
            {record?.currency} {parseFloat(record?.total_disbursed || 0).toLocaleString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="Members">
        <AntList
          dataSource={record?.members || []}
          renderItem={(member: any) => (
            <AntList.Item>
              <AntList.Item.Meta
                title={`Position #${member.payout_position}: ${member.user?.profile?.first_name} ${member.user?.profile?.last_name}`}
                description={`Status: ${member.status} | Contributed: ${record?.currency} ${parseFloat(member.total_contributed || 0).toLocaleString()}`}
              />
              {member.payout_received && (
                <Tag color="green">Payout Received</Tag>
              )}
            </AntList.Item>
          )}
        />
      </Card>
    </Show>
  );
};
