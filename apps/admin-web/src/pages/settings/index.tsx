import { Card, Form, Input, Switch, Button, Tabs, message } from "antd";

export const Settings: React.FC = () => {
  const [generalForm] = Form.useForm();
  const [fxForm] = Form.useForm();

  const handleSaveGeneral = (values: any) => {
    console.log("General settings:", values);
    message.success("General settings saved");
  };

  const handleSaveFX = (values: any) => {
    console.log("FX settings:", values);
    message.success("FX settings saved");
  };

  return (
    <div style={{ padding: 24 }}>
      <Tabs
        items={[
          {
            key: "general",
            label: "General",
            children: (
              <Card title="General Settings">
                <Form
                  form={generalForm}
                  layout="vertical"
                  onFinish={handleSaveGeneral}
                  initialValues={{
                    platformName: "HustleX",
                    supportEmail: "support@hustlex.ng",
                    supportPhone: "+234 700 HUSTLEX",
                    maintenanceMode: false,
                  }}
                >
                  <Form.Item label="Platform Name" name="platformName">
                    <Input />
                  </Form.Item>
                  <Form.Item label="Support Email" name="supportEmail">
                    <Input />
                  </Form.Item>
                  <Form.Item label="Support Phone" name="supportPhone">
                    <Input />
                  </Form.Item>
                  <Form.Item
                    label="Maintenance Mode"
                    name="maintenanceMode"
                    valuePropName="checked"
                  >
                    <Switch />
                  </Form.Item>
                  <Button type="primary" htmlType="submit">
                    Save Changes
                  </Button>
                </Form>
              </Card>
            ),
          },
          {
            key: "fx",
            label: "FX & Remittance",
            children: (
              <Card title="FX Settings">
                <Form
                  form={fxForm}
                  layout="vertical"
                  onFinish={handleSaveFX}
                  initialValues={{
                    gbpNgnSpread: 150,
                    usdNgnSpread: 175,
                    eurNgnSpread: 175,
                    cadNgnSpread: 200,
                    quoteValidityMinutes: 15,
                    minTransferAmount: 10,
                    maxTransferAmount: 10000,
                    transferFeeFixed: 2.99,
                    transferFeePercent: 0.5,
                  }}
                >
                  <Form.Item label="GBP to NGN Spread (bps)" name="gbpNgnSpread">
                    <Input type="number" />
                  </Form.Item>
                  <Form.Item label="USD to NGN Spread (bps)" name="usdNgnSpread">
                    <Input type="number" />
                  </Form.Item>
                  <Form.Item label="EUR to NGN Spread (bps)" name="eurNgnSpread">
                    <Input type="number" />
                  </Form.Item>
                  <Form.Item label="CAD to NGN Spread (bps)" name="cadNgnSpread">
                    <Input type="number" />
                  </Form.Item>
                  <Form.Item label="Quote Validity (minutes)" name="quoteValidityMinutes">
                    <Input type="number" />
                  </Form.Item>
                  <Form.Item label="Min Transfer Amount" name="minTransferAmount">
                    <Input type="number" />
                  </Form.Item>
                  <Form.Item label="Max Transfer Amount" name="maxTransferAmount">
                    <Input type="number" />
                  </Form.Item>
                  <Form.Item label="Fixed Transfer Fee" name="transferFeeFixed">
                    <Input type="number" step="0.01" />
                  </Form.Item>
                  <Form.Item label="Transfer Fee %" name="transferFeePercent">
                    <Input type="number" step="0.1" />
                  </Form.Item>
                  <Button type="primary" htmlType="submit">
                    Save FX Settings
                  </Button>
                </Form>
              </Card>
            ),
          },
          {
            key: "notifications",
            label: "Notifications",
            children: (
              <Card title="Notification Settings">
                <Form layout="vertical">
                  <Form.Item
                    label="Enable SMS Notifications"
                    name="smsEnabled"
                    valuePropName="checked"
                    initialValue={true}
                  >
                    <Switch />
                  </Form.Item>
                  <Form.Item
                    label="Enable Push Notifications"
                    name="pushEnabled"
                    valuePropName="checked"
                    initialValue={true}
                  >
                    <Switch />
                  </Form.Item>
                  <Form.Item
                    label="Enable Email Notifications"
                    name="emailEnabled"
                    valuePropName="checked"
                    initialValue={true}
                  >
                    <Switch />
                  </Form.Item>
                  <Form.Item label="SMS Provider" name="smsProvider">
                    <Input placeholder="africa_talking" />
                  </Form.Item>
                  <Button type="primary" htmlType="submit">
                    Save Notification Settings
                  </Button>
                </Form>
              </Card>
            ),
          },
        ]}
      />
    </div>
  );
};
