import { Refine } from "@refinedev/core";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";
import {
  ErrorComponent,
  ThemedLayoutV2,
  useNotificationProvider,
} from "@refinedev/antd";
import { App as AntdApp, ConfigProvider } from "antd";
import routerBindings, {
  DocumentTitleHandler,
  UnsavedChangesNotifier,
} from "@refinedev/react-router-v6";
import { BrowserRouter, Outlet, Route, Routes } from "react-router-dom";
import {
  DashboardOutlined,
  UserOutlined,
  ShopOutlined,
  DollarOutlined,
  CalendarOutlined,
  BankOutlined,
  SendOutlined,
  TeamOutlined,
  BellOutlined,
  SettingOutlined,
} from "@ant-design/icons";

import "@refinedev/antd/dist/reset.css";

import { dataProvider, liveProvider } from "./providers/hasuraProvider";
import { authProvider } from "./providers/authProvider";

// Pages
import { Dashboard } from "./pages/dashboard";
import { UserList, UserShow, UserEdit } from "./pages/users";
import { ServiceList, ServiceShow, ServiceEdit } from "./pages/services";
import { BookingList, BookingShow } from "./pages/bookings";
import { TransactionList, TransactionShow } from "./pages/transactions";
import { RemittanceList, RemittanceShow } from "./pages/remittances";
import { BeneficiaryList, BeneficiaryShow } from "./pages/beneficiaries";
import { SavingsCircleList, SavingsCircleShow } from "./pages/savingsCircles";
import { NotificationList } from "./pages/notifications";
import { Settings } from "./pages/settings";
import { Login } from "./pages/auth/login";

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <RefineKbarProvider>
        <ConfigProvider
          theme={{
            token: {
              colorPrimary: "#4F46E5",
              borderRadius: 8,
            },
          }}
        >
          <AntdApp>
            <Refine
              dataProvider={dataProvider}
              liveProvider={liveProvider}
              authProvider={authProvider}
              notificationProvider={useNotificationProvider}
              routerProvider={routerBindings}
              resources={[
                {
                  name: "dashboard",
                  list: "/",
                  meta: {
                    label: "Dashboard",
                    icon: <DashboardOutlined />,
                  },
                },
                {
                  name: "users",
                  list: "/users",
                  show: "/users/:id",
                  edit: "/users/:id/edit",
                  meta: {
                    label: "Users",
                    icon: <UserOutlined />,
                  },
                },
                {
                  name: "services",
                  list: "/services",
                  show: "/services/:id",
                  edit: "/services/:id/edit",
                  meta: {
                    label: "Services",
                    icon: <ShopOutlined />,
                  },
                },
                {
                  name: "bookings",
                  list: "/bookings",
                  show: "/bookings/:id",
                  meta: {
                    label: "Bookings",
                    icon: <CalendarOutlined />,
                  },
                },
                {
                  name: "transactions",
                  list: "/transactions",
                  show: "/transactions/:id",
                  meta: {
                    label: "Transactions",
                    icon: <DollarOutlined />,
                  },
                },
                {
                  name: "remittances",
                  list: "/remittances",
                  show: "/remittances/:id",
                  meta: {
                    label: "Remittances",
                    icon: <SendOutlined />,
                  },
                },
                {
                  name: "beneficiaries",
                  list: "/beneficiaries",
                  show: "/beneficiaries/:id",
                  meta: {
                    label: "Beneficiaries",
                    icon: <BankOutlined />,
                  },
                },
                {
                  name: "savings_circles",
                  list: "/savings-circles",
                  show: "/savings-circles/:id",
                  meta: {
                    label: "Savings Circles",
                    icon: <TeamOutlined />,
                  },
                },
                {
                  name: "notifications",
                  list: "/notifications",
                  meta: {
                    label: "Notifications",
                    icon: <BellOutlined />,
                  },
                },
                {
                  name: "settings",
                  list: "/settings",
                  meta: {
                    label: "Settings",
                    icon: <SettingOutlined />,
                  },
                },
              ]}
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
                liveMode: "auto",
              }}
            >
              <Routes>
                <Route
                  element={
                    <ThemedLayoutV2
                      Title={() => (
                        <div style={{ fontSize: 18, fontWeight: 700 }}>
                          HustleX Admin
                        </div>
                      )}
                    >
                      <Outlet />
                    </ThemedLayoutV2>
                  }
                >
                  <Route index element={<Dashboard />} />
                  <Route path="/users">
                    <Route index element={<UserList />} />
                    <Route path=":id" element={<UserShow />} />
                    <Route path=":id/edit" element={<UserEdit />} />
                  </Route>
                  <Route path="/services">
                    <Route index element={<ServiceList />} />
                    <Route path=":id" element={<ServiceShow />} />
                    <Route path=":id/edit" element={<ServiceEdit />} />
                  </Route>
                  <Route path="/bookings">
                    <Route index element={<BookingList />} />
                    <Route path=":id" element={<BookingShow />} />
                  </Route>
                  <Route path="/transactions">
                    <Route index element={<TransactionList />} />
                    <Route path=":id" element={<TransactionShow />} />
                  </Route>
                  <Route path="/remittances">
                    <Route index element={<RemittanceList />} />
                    <Route path=":id" element={<RemittanceShow />} />
                  </Route>
                  <Route path="/beneficiaries">
                    <Route index element={<BeneficiaryList />} />
                    <Route path=":id" element={<BeneficiaryShow />} />
                  </Route>
                  <Route path="/savings-circles">
                    <Route index element={<SavingsCircleList />} />
                    <Route path=":id" element={<SavingsCircleShow />} />
                  </Route>
                  <Route path="/notifications" element={<NotificationList />} />
                  <Route path="/settings" element={<Settings />} />
                  <Route path="*" element={<ErrorComponent />} />
                </Route>
                <Route path="/login" element={<Login />} />
              </Routes>
              <RefineKbar />
              <UnsavedChangesNotifier />
              <DocumentTitleHandler />
            </Refine>
          </AntdApp>
        </ConfigProvider>
      </RefineKbarProvider>
    </BrowserRouter>
  );
};

export default App;
