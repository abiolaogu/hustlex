import { Row, Col, Card, Statistic, Typography, Table, Tag } from "antd";
import {
  UserOutlined,
  DollarOutlined,
  ShopOutlined,
  CalendarOutlined,
  ArrowUpOutlined,
} from "@ant-design/icons";

const { Title } = Typography;

export const Dashboard: React.FC = () => {
  // Mock data - will be replaced with actual API calls
  const stats = {
    totalUsers: 12450,
    activeProviders: 3420,
    totalTransactions: 458900,
    monthlyRevenue: 45600000,
    userGrowth: 12.5,
    revenueGrowth: 8.3,
  };

  const recentBookings = [
    { id: "1", service: "House Cleaning", status: "confirmed", amount: 15000 },
    { id: "2", service: "Hair Styling", status: "completed", amount: 8000 },
    { id: "3", service: "Plumbing", status: "pending", amount: 25000 },
    { id: "4", service: "Photography", status: "in_progress", amount: 50000 },
  ];

  const recentRemittances = [
    { id: "1", sender: "John D.", amount: "£500", target: "₦975,000", status: "completed" },
    { id: "2", sender: "Mary O.", amount: "$300", target: "₦465,000", status: "processing" },
    { id: "3", sender: "Peter A.", amount: "€400", target: "₦680,000", status: "pending" },
  ];

  const statusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: "orange",
      confirmed: "blue",
      in_progress: "cyan",
      completed: "green",
      cancelled: "red",
      processing: "purple",
    };
    return colors[status] || "default";
  };

  return (
    <div style={{ padding: 24 }}>
      <Title level={2}>Dashboard</Title>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Total Users"
              value={stats.totalUsers}
              prefix={<UserOutlined />}
              suffix={
                <span style={{ fontSize: 14, color: "#52c41a" }}>
                  <ArrowUpOutlined /> {stats.userGrowth}%
                </span>
              }
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Active Providers"
              value={stats.activeProviders}
              prefix={<ShopOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Total Transactions"
              value={stats.totalTransactions}
              prefix={<CalendarOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Monthly Revenue"
              value={stats.monthlyRevenue}
              prefix={<DollarOutlined />}
              precision={0}
              suffix={
                <span style={{ fontSize: 14, color: "#52c41a" }}>
                  <ArrowUpOutlined /> {stats.revenueGrowth}%
                </span>
              }
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} lg={12}>
          <Card title="Recent Bookings">
            <Table
              dataSource={recentBookings}
              columns={[
                { title: "Service", dataIndex: "service", key: "service" },
                {
                  title: "Status",
                  dataIndex: "status",
                  key: "status",
                  render: (status) => (
                    <Tag color={statusColor(status)}>
                      {status.replace("_", " ").toUpperCase()}
                    </Tag>
                  ),
                },
                {
                  title: "Amount",
                  dataIndex: "amount",
                  key: "amount",
                  render: (amount) => `₦${amount.toLocaleString()}`,
                },
              ]}
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="Recent Remittances">
            <Table
              dataSource={recentRemittances}
              columns={[
                { title: "Sender", dataIndex: "sender", key: "sender" },
                { title: "Sent", dataIndex: "amount", key: "amount" },
                { title: "Received", dataIndex: "target", key: "target" },
                {
                  title: "Status",
                  dataIndex: "status",
                  key: "status",
                  render: (status) => (
                    <Tag color={statusColor(status)}>
                      {status.toUpperCase()}
                    </Tag>
                  ),
                },
              ]}
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
