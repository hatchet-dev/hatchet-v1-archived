import { AppWrapper } from "components/appwrapper";
import React, { Component } from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import { GlobalStyle } from "shared/theme";
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

export default class App extends Component {
  state = {
    currentSection: "Modules",
  };

  renderHomeContents = () => {
    return (
      <>
        <Switch>
          <Route path="/monitoring" render={() => <MonitoringView />}></Route>
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
    );
  };

  render() {
    return (
      <AppWrapper>
        <BrowserRouter>
          <TopBar />
          <SideBar />
          <GlobalStyle />
          <ViewWrapper>{this.renderHomeContents()}</ViewWrapper>
        </BrowserRouter>
      </AppWrapper>
    );
  }
}
