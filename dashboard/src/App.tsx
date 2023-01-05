import { AppWrapper } from "components/appwrapper";
import React, { Component } from "react";
import {
  BrowserRouter,
  Route,
  Switch,
  useHistory,
  useLocation,
} from "react-router-dom";
import theme, { GlobalStyle } from "shared/theme";
import TemplatesView from "views/templates/TemplatesView";
import ModulesView from "views/modules/ModulesView";
import TopBar from "components/topbar";
import SideBar, { SidebarLink } from "components/sidebar";
import { ViewWrapper } from "components/viewwrapper";
import EnvironmentsView from "views/environments/EnvironmentsView";
import HomeView from "views/home/HomeView";
import MonitoringView from "views/monitoring/MonitoringView";
import LinkModuleView from "views/linkmodule/LinkModuleView";
import ExpandedModuleView from "views/expandedmodule/ExpandedModuleView";
import ExpandedTemplateView from "views/expandedtemplate/ExpandedTemplateView";
import LoginView from "views/login/LoginView";
import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from "@tanstack/react-query";
import RegisterView from "views/register/RegisterView";
import AuthChecker from "shared/auth/AuthChecker";
import api from "shared/api";
import { ThemeProvider } from "styled-components";
import UserSettingsView from "views/usersettings/UserSettingsView";

const App: React.FunctionComponent = () => {
  const queryClient = new QueryClient();

  return (
    <ThemeProvider theme={theme}>
      <AppWrapper>
        <GlobalStyle />
        <BrowserRouter>
          <QueryClientProvider client={queryClient}>
            <AppContents />
          </QueryClientProvider>
        </BrowserRouter>
      </AppWrapper>
    </ThemeProvider>
  );
};

export default App;

const DashboardSidebarLinks: SidebarLink[] = [
  {
    name: "Home",
    href: "/home",
  },
  {
    name: "Modules",
    href: "/modules",
  },
  {
    name: "Monitoring",
    href: "/monitoring",
  },
  {
    name: "Templates",
    href: "/templates",
  },
  {
    name: "Integrations",
    href: "/integrations",
  },
  {
    name: "Settings",
    href: "/settings",
  },
];

const UserSidebarLinks: SidebarLink[] = [
  {
    name: "Profile",
    href: "/user/settings",
  },
  {
    name: "Personal Access Tokens",
    href: "/user/settings/pats",
  },
];

const AppContents: React.FunctionComponent = () => {
  const renderAppContents = () => {
    return (
      <>
        <Switch>
          <Route
            path="/login"
            render={() => (
              <AuthChecker check_authenticated={false}>
                <LoginView />
              </AuthChecker>
            )}
          ></Route>
          <Route
            path="/register"
            render={() => (
              <AuthChecker check_authenticated={false}>
                <RegisterView />
              </AuthChecker>
            )}
          ></Route>
          <Route
            path="/user/*"
            render={() => renderUserSettingsContents()}
          ></Route>
          <Route path="/" render={() => renderHomeContents()}></Route>
        </Switch>
      </>
    );
  };

  const renderUserSettingsContents = () => {
    return (
      <AuthChecker check_authenticated={true}>
        <TopBar />
        <SideBar links={UserSidebarLinks} />
        <ViewWrapper>
          <>
            <Switch>
              <Route
                path="/user/settings"
                render={() => <UserSettingsView />}
              ></Route>
            </Switch>
          </>
        </ViewWrapper>
      </AuthChecker>
    );
  };

  const renderHomeContents = () => {
    return (
      <AuthChecker check_authenticated={true}>
        <TopBar />
        <SideBar links={DashboardSidebarLinks} />
        <ViewWrapper>
          <>
            <Switch>
              <Route
                path="/monitoring"
                render={() => <MonitoringView />}
              ></Route>
              <Route
                path="/modules/link/:step"
                render={() => <LinkModuleView />}
              ></Route>
              <Route
                path="/modules/:module"
                render={() => <ExpandedModuleView />}
              ></Route>
              <Route path="/modules" render={() => <ModulesView />}></Route>
              <Route
                path="/templates/:template"
                render={() => <ExpandedTemplateView />}
              ></Route>
              <Route path="/templates" render={() => <TemplatesView />}></Route>
              <Route
                path="/environments"
                render={() => <EnvironmentsView />}
              ></Route>
              <Route path="/home" render={() => <HomeView />}></Route>
              <Route path="/" render={() => <HomeView />}></Route>
            </Switch>
          </>
        </ViewWrapper>
      </AuthChecker>
    );
  };

  return renderAppContents();
};
