import {
  AppWrapper,
  ViewWrapper,
  TopBar,
} from "@hatchet-dev/hatchet-components";
import React from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import theme, { GlobalStyle } from "shared/theme";
import TemplatesView from "views/templates/TemplatesView";
import ModulesView from "views/modules/ModulesView";
import TopBarWithProfile from "components/topbarwithprofile";
import SideBar, { SidebarLink } from "components/sidebar";
import EnvironmentsView from "views/environments/EnvironmentsView";
import HomeView from "views/home/HomeView";
import MonitoringView from "views/monitoring/MonitoringView";
import LinkModuleView from "views/linkmodule/LinkModuleView";
import ExpandedModuleView from "views/expandedmodule/ExpandedModuleView";
import ExpandedTemplateView from "views/expandedtemplate/ExpandedTemplateView";
import LoginView from "views/login/LoginView";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import RegisterView from "views/register/RegisterView";
import AuthChecker from "shared/auth/AuthChecker";
import { ThemeProvider } from "styled-components";
import UserSettingsView from "views/usersettings/UserSettingsView";
import CreateOrganizationView from "views/createorganization/CreateOrganizationView";
import Populator from "shared/populator/Populator";
import UserPATsView from "views/userpatsview/UserPATsView";
import InitiateResetPasswordView from "views/initiateresetpassword/InitiateResetPasswordView";
import FinalizeResetPasswordView from "views/finalizeresetpassword/FinalizeResetPasswordView";
import VerifyEmailView from "views/verifyemail/VerifyEmail";
import AcceptOrganizationInviteView from "views/acceptorganizationinvite/AcceptOrganizationInviteView";
import OrganizationSettingsView from "views/organizationsettings/OrganizationSettingsView";
import SettingsView from "views/settings/SettingsView";

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
    href: "/user/settings/profile",
  },
  {
    name: "Linked Accounts",
    href: "/user/settings/accounts",
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
                <TopBar />
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
            path="/reset_password/initiate"
            render={() => (
              <AuthChecker check_authenticated={false}>
                <InitiateResetPasswordView />
              </AuthChecker>
            )}
          ></Route>
          <Route
            path="/reset_password/finalize"
            render={() => (
              <AuthChecker check_authenticated={false}>
                <FinalizeResetPasswordView />
              </AuthChecker>
            )}
          ></Route>
          <Route
            path="/verify_email/finalize"
            render={() => (
              <AuthChecker
                check_authenticated={true}
                allow_unverified_email={true}
                require_unverified_email={true}
              >
                <VerifyEmailView />
              </AuthChecker>
            )}
          ></Route>
          <Route
            path="/organization_invite/accept"
            render={() => (
              <AuthChecker check_authenticated={true}>
                <AcceptOrganizationInviteView />
              </AuthChecker>
            )}
          ></Route>
          <Route
            path="/user/*"
            render={() => renderUserSettingsContents()}
          ></Route>
          <Route
            path="/organization/create"
            render={() => renderOnboardingContents()}
          ></Route>
          <Route path="/" render={() => renderHomeContents()}></Route>
        </Switch>
      </>
    );
  };

  const renderUserSettingsContents = () => {
    return (
      <AuthChecker check_authenticated={true}>
        <TopBarWithProfile />
        <SideBar links={UserSidebarLinks} />
        <ViewWrapper>
          <>
            <Switch>
              <Route
                path="/user/settings/pats"
                render={() => <UserPATsView />}
              ></Route>
              <Route
                path="/user/settings/profile"
                render={() => <UserSettingsView />}
              ></Route>
            </Switch>
          </>
        </ViewWrapper>
      </AuthChecker>
    );
  };

  const renderOnboardingContents = () => {
    return (
      <AuthChecker check_authenticated={true}>
        <TopBarWithProfile />
        <ViewWrapper>
          <>
            <Switch>
              <Route
                path="/organization/create"
                render={() => <CreateOrganizationView />}
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
        <TopBarWithProfile />
        <Populator organization team>
          <SideBar links={DashboardSidebarLinks} />
          <ViewWrapper>
            <>
              <Switch>
                <Route
                  path="/settings/organization"
                  render={() => (
                    <SettingsView defaultOption="Organization Settings" />
                  )}
                ></Route>
                <Route
                  path="/settings/team"
                  render={() => <SettingsView defaultOption="Team Settings" />}
                ></Route>
                <Route
                  path="/settings"
                  render={() => (
                    <SettingsView defaultOption="Organization Settings" />
                  )}
                ></Route>
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
                <Route
                  path="/templates"
                  render={() => <TemplatesView />}
                ></Route>
                <Route
                  path="/environments"
                  render={() => <EnvironmentsView />}
                ></Route>
                <Route path="/home" render={() => <HomeView />}></Route>
                <Route path="/" render={() => <HomeView />}></Route>
              </Switch>
            </>
          </ViewWrapper>
        </Populator>
      </AuthChecker>
    );
  };

  return renderAppContents();
};
