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
import CreateModuleView from "views/createmodule/CreateModuleView";
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
import TeamSettingsView from "views/teamsettings/TeamSettingsView";
import TeamsList from "views/teamslist/TeamsList";
import CreateTeam from "views/createteam/CreateTeam";
import UserAccountsView from "views/useraccounts/UserAccountsView";
import CreateMonitorView from "views/createmonitor/CreateMonitorView";
import ExpandedMonitorView from "views/expandedmonitor/ExpandedMonitorView";
import Notifications from "views/notifications/Notifications";
import HatchetErrorBoundary from "shared/errors/ErrorBoundary";

const App: React.FunctionComponent = () => {
  const queryClient = new QueryClient();

  return (
    <ThemeProvider theme={theme}>
      <AppWrapper>
        <GlobalStyle />
        <BrowserRouter>
          <QueryClientProvider client={queryClient}>
            <HatchetErrorBoundary>
              <AppContents />
            </HatchetErrorBoundary>
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
    name: "Notifications",
    href: "/notifications",
  },
  {
    name: "Teams",
    href: "/organization/teams",
  },
  {
    name: "Organization Settings",
    href: "/organization/settings",
  },
];

const DashboardTeamSidebarLinks: SidebarLink[] = [
  // {
  //   name: "Overview",
  //   href: "/overview",
  // },
  {
    name: "Modules",
    href: "/modules",
  },
  {
    name: "Monitors",
    href: "/monitors",
  },
  // {
  //   name: "Templates",
  //   href: "/templates",
  // },
  // {
  //   name: "Integrations",
  //   href: "/integrations",
  // },
  {
    name: "Team Settings",
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
                path="/user/settings/accounts"
                render={() => <UserAccountsView />}
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
        <Switch>
          <Route
            path="/organization/settings"
            render={() => (
              <WrappedTeamContents>
                <OrganizationSettingsView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/organization/teams/create"
            render={() => (
              <WrappedTeamContents>
                <CreateTeam />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/organization/teams"
            render={() => (
              <WrappedTeamContents>
                <TeamsList />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/notifications/teams/:team/:notification"
            render={() => (
              <WrappedTeamContents>
                <Notifications />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/notifications"
            render={() => (
              <WrappedTeamContents>
                <Notifications />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/settings"
            render={() => (
              <WrappedTeamContents>
                <TeamSettingsView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/monitors/create/:step"
            render={() => (
              <WrappedTeamContents>
                <CreateMonitorView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/monitors/:monitor"
            render={() => (
              <WrappedTeamContents>
                <ExpandedMonitorView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/monitors"
            render={() => (
              <WrappedTeamContents>
                <MonitoringView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/modules/create/:step"
            render={() => (
              <WrappedTeamContents>
                <CreateModuleView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/modules/:module"
            render={() => (
              <WrappedTeamContents>
                <ExpandedModuleView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/modules"
            render={() => (
              <WrappedTeamContents>
                <ModulesView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/templates/:template"
            render={() => (
              <WrappedTeamContents>
                <ExpandedTemplateView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/templates"
            render={() => (
              <WrappedTeamContents>
                <TemplatesView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team/environments"
            render={() => (
              <WrappedTeamContents>
                <EnvironmentsView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/teams/:team"
            render={() => (
              <WrappedTeamContents>
                <ModulesView />
              </WrappedTeamContents>
            )}
          ></Route>
          <Route
            path="/"
            render={() => (
              <WrappedTeamContents>
                <HomeView />
              </WrappedTeamContents>
            )}
          ></Route>
        </Switch>
      </AuthChecker>
    );
  };

  return renderAppContents();
};

type WrapperContentsProps = {
  children?: React.ReactNode;
};

const WrappedTeamContents: React.FunctionComponent<WrapperContentsProps> = ({
  children,
}) => {
  return (
    <>
      <Populator organization team>
        <SideBar
          links={DashboardSidebarLinks}
          team_links={DashboardTeamSidebarLinks}
        />
        <HatchetErrorBoundary>
          <ViewWrapper>{children}</ViewWrapper>
        </HatchetErrorBoundary>
      </Populator>
    </>
  );
};
