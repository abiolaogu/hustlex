import { List, useTable } from "@refinedev/antd";
import { Table, Tag } from "antd";

export const NotificationList: React.FC = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  const statusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: "orange",
      sent: "blue",
      delivered: "green",
      failed: "red",
      read: "cyan",
    };
    return colors[status] || "default";
  };

  const typeColor = (type: string) => {
    const colors: Record<string, string> = {
      push: "purple",
      sms: "blue",
      email: "cyan",
      in_app: "green",
    };
    return colors[type] || "default";
  };

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column
          dataIndex={["user", "profile", "first_name"]}
          title="User"
          render={(_, record: any) =>
            `${record.user?.profile?.first_name || ""} ${record.user?.profile?.last_name || ""}`
          }
        />
        <Table.Column
          dataIndex="type"
          title="Type"
          render={(type) => <Tag color={typeColor(type)}>{type?.toUpperCase()}</Tag>}
        />
        <Table.Column dataIndex="title" title="Title" />
        <Table.Column
          dataIndex="body"
          title="Body"
          ellipsis
          width={300}
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
          title="Created"
          render={(date) => new Date(date).toLocaleString()}
        />
        <Table.Column
          dataIndex="sent_at"
          title="Sent"
          render={(date) => date ? new Date(date).toLocaleString() : "-"}
        />
      </Table>
    </List>
  );
};
