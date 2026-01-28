import { useLogin } from "@refinedev/core";
import { Form, Input, Button, Card, Typography, message } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";

const { Title } = Typography;

export const Login: React.FC = () => {
  const { mutate: login, isLoading } = useLogin();

  const onFinish = (values: { email: string; password: string }) => {
    login(values, {
      onError: () => {
        message.error("Login failed. Please check your credentials.");
      },
    });
  };

  return (
    <div
      style={{
        height: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
      }}
    >
      <Card style={{ width: 400, padding: 24 }}>
        <div style={{ textAlign: "center", marginBottom: 24 }}>
          <Title level={2} style={{ margin: 0 }}>
            HustleX Admin
          </Title>
          <Typography.Text type="secondary">
            Sign in to continue
          </Typography.Text>
        </div>

        <Form
          name="login"
          onFinish={onFinish}
          layout="vertical"
          initialValues={{
            email: "admin@hustlex.ng",
            password: "admin123",
          }}
        >
          <Form.Item
            name="email"
            rules={[
              { required: true, message: "Please enter your email" },
              { type: "email", message: "Please enter a valid email" },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="Email"
              size="large"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: "Please enter your password" }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="Password"
              size="large"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              size="large"
              block
              loading={isLoading}
            >
              Sign In
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};
