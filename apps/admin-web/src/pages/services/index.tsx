import {
  List,
  useTable,
  EditButton,
  ShowButton,
  Show,
  Edit,
  useForm,
} from "@refinedev/antd";
import { Table, Space, Tag, Form, Input, Select, Card, Descriptions, Rate } from "antd";
import { useShow } from "@refinedev/core";

export const ServiceList: React.FC = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  const statusColor = (status: string) => {
    const colors: Record<string, string> = {
      active: "green",
      draft: "default",
      paused: "orange",
      archived: "red",
    };
    return colors[status] || "default";
  };

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="title" title="Title" />
        <Table.Column
          dataIndex={["category", "name"]}
          title="Category"
        />
        <Table.Column
          dataIndex={["provider", "profile", "first_name"]}
          title="Provider"
          render={(_, record: any) =>
            `${record.provider?.profile?.first_name || ""} ${record.provider?.profile?.last_name || ""}`
          }
        />
        <Table.Column
          dataIndex="base_price"
          title="Price"
          render={(price, record: any) =>
            `${record.currency || "NGN"} ${parseFloat(price).toLocaleString()}`
          }
        />
        <Table.Column
          dataIndex="average_rating"
          title="Rating"
          render={(rating) => <Rate disabled defaultValue={rating} allowHalf />}
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
              <EditButton hideText size="small" recordItemId={record.id} />
              <ShowButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};

export const ServiceShow: React.FC = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Card title="Service Details">
        <Descriptions bordered column={2}>
          <Descriptions.Item label="Title">{record?.title}</Descriptions.Item>
          <Descriptions.Item label="Category">{record?.category?.name}</Descriptions.Item>
          <Descriptions.Item label="Price">
            {record?.currency} {parseFloat(record?.base_price || 0).toLocaleString()}
          </Descriptions.Item>
          <Descriptions.Item label="Rating">
            <Rate disabled defaultValue={record?.average_rating} allowHalf />
          </Descriptions.Item>
          <Descriptions.Item label="Status">{record?.status}</Descriptions.Item>
          <Descriptions.Item label="Bookings">{record?.booking_count}</Descriptions.Item>
          <Descriptions.Item label="Description" span={2}>
            {record?.description}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </Show>
  );
};

export const ServiceEdit: React.FC = () => {
  const { formProps, saveButtonProps } = useForm();

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item label="Title" name="title" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item label="Description" name="description">
          <Input.TextArea rows={4} />
        </Form.Item>
        <Form.Item label="Base Price" name="base_price">
          <Input type="number" />
        </Form.Item>
        <Form.Item label="Status" name="status">
          <Select
            options={[
              { label: "Active", value: "active" },
              { label: "Draft", value: "draft" },
              { label: "Paused", value: "paused" },
              { label: "Archived", value: "archived" },
            ]}
          />
        </Form.Item>
      </Form>
    </Edit>
  );
};
