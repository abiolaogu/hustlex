import {
  List,
  useTable,
  EditButton,
  ShowButton,
  DeleteButton,
  FilterDropdown,
  useSelect,
  Show,
  Edit,
  useForm,
} from "@refinedev/antd";
import { Table, Space, Tag, Input, Select, Form, Card, Descriptions } from "antd";
import { useShow, useOne } from "@refinedev/core";

export const UserList: React.FC = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  const statusColor = (status: string) => {
    const colors: Record<string, string> = {
      active: "green",
      inactive: "default",
      suspended: "red",
      pending_verification: "orange",
    };
    return colors[status] || "default";
  };

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="ID" />
        <Table.Column
          dataIndex={["profile", "first_name"]}
          title="Name"
          render={(_, record: any) =>
            `${record.profile?.first_name || ""} ${record.profile?.last_name || ""}`
          }
        />
        <Table.Column dataIndex="email" title="Email" />
        <Table.Column dataIndex="phone" title="Phone" />
        <Table.Column
          dataIndex="role"
          title="Role"
          render={(role) => <Tag>{role?.toUpperCase()}</Tag>}
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
              <EditButton hideText size="small" recordItemId={record.id} />
              <ShowButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};

export const UserShow: React.FC = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Card title="User Information">
        <Descriptions bordered column={2}>
          <Descriptions.Item label="ID">{record?.id}</Descriptions.Item>
          <Descriptions.Item label="Email">{record?.email}</Descriptions.Item>
          <Descriptions.Item label="Phone">{record?.phone}</Descriptions.Item>
          <Descriptions.Item label="Role">{record?.role}</Descriptions.Item>
          <Descriptions.Item label="Status">{record?.status}</Descriptions.Item>
          <Descriptions.Item label="Created">
            {new Date(record?.created_at).toLocaleDateString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </Show>
  );
};

export const UserEdit: React.FC = () => {
  const { formProps, saveButtonProps } = useForm();

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item label="Status" name="status">
          <Select
            options={[
              { label: "Active", value: "active" },
              { label: "Inactive", value: "inactive" },
              { label: "Suspended", value: "suspended" },
              { label: "Pending Verification", value: "pending_verification" },
            ]}
          />
        </Form.Item>
        <Form.Item label="Role" name="role">
          <Select
            options={[
              { label: "Consumer", value: "consumer" },
              { label: "Provider", value: "provider" },
              { label: "Admin", value: "admin" },
            ]}
          />
        </Form.Item>
      </Form>
    </Edit>
  );
};
