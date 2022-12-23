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
import SideBar from "components/sidebar";
import { ViewWrapper } from "components/viewwrapper";
import EnvironmentsView from "views/environments/EnvironmentsView";
import HomeView from "views/homeview/HomeView";
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

const AppContents: React.FunctionComponent = () => {
  const history = useHistory();
  const location = useLocation();

  const authCheckEnabled =
    location.pathname != "/login" && location.pathname != "/register";

  const { status, data, error, isFetching } = useQuery({
    queryKey: ["initial_current_user"],
    queryFn: async () => {
      const res = await api.getCurrentUser();
      return res;
    },
    retry: false,
    enabled: authCheckEnabled,
  });

  // TODO(abelanger5): style/case on loading
  if (isFetching) {
    console.log("IS FETCHING");
    return <div>Loading...</div>;
  }

  if (authCheckEnabled && error) {
    history.push("/login");
  }

  const renderAppContents = () => {
    return (
      <>
        <Switch>
          <Route path="/login" render={() => <LoginView />}></Route>
          <Route path="/register" render={() => <RegisterView />}></Route>
          <Route path="/" render={() => renderHomeContents()}></Route>
        </Switch>
      </>
    );
  };

  const renderHomeContents = () => {
    return (
      <AuthChecker>
        <TopBar />
        <SideBar />
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
