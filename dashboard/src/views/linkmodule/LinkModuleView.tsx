import Breadcrumbs from "components/breadcrumbs";
import { FlexRowRight, H1, HorizontalSpacer } from "components/globals";
import React from "react";
import { useHistory, useParams } from "react-router-dom";
import StandardButton from "components/buttons";
import ChooseGitSource from "./components/ChooseGitSource";
import LinkVariables from "./components/LinkVariables";

const LinkModuleView: React.FunctionComponent = () => {
  let history = useHistory();
  const { step } = useParams<{ step: string }>();

  const renderForm = () => {
    switch (step) {
      case "step_1":
        return <ChooseGitSource />;
      case "step_2":
        return <LinkVariables />;
    }
  };

  const getBreadcrumbs = () => {
    switch (step) {
      case "step_1":
        return [
          {
            label: "Modules",
            link: "/modules",
          },
          {
            label: "Step 1: Choose Git Source",
            link: "/modules/link/step_1",
          },
        ];
      case "step_2":
        return [
          {
            label: "Modules",
            link: "/modules",
          },
          {
            label: "Step 1: Choose Git Source",
            link: "/modules/link/step_1",
          },
          {
            label: "Step 2: Choose Variable Source",
            link: "/modules/link/step_1",
          },
        ];
    }
  };

  const onSubmit = () => {
    history.push("/modules");
  };

  const renderButton = () => {
    switch (step) {
      case "step_1":
        return (
          <StandardButton
            label="Next"
            material_icon="chevron_right"
            icon_side="right"
            on_click={() => {
              // TODO: store parameters from step 1
              history.push("/modules/link/step_2");
            }}
          />
        );
      case "step_2":
        return (
          <StandardButton
            label="Submit"
            on_click={() => {
              // TODO: save form
              history.push("/modules");
            }}
          />
        );
    }
  };

  return (
    <>
      <Breadcrumbs breadcrumbs={getBreadcrumbs()} />
      <HorizontalSpacer spacepixels={12} />
      <H1>Link a new module</H1>
      <HorizontalSpacer spacepixels={20} />
      {renderForm()}
      <HorizontalSpacer spacepixels={20} />
      <FlexRowRight>{renderButton()}</FlexRowRight>
    </>
  );
};

export default LinkModuleView;
