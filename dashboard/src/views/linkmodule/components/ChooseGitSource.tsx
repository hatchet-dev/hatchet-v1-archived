import { H2, HorizontalSpacer, P } from "components/globals";
import Selector from "components/selector";
import React from "react";
import { useHistory } from "react-router-dom";
import gitRepository from "assets/git_repository.png";
import github from "assets/github.png";
import branch from "assets/branch.png";
import TextInput from "components/textinput";
import { css } from "styled-components";
import theme from "shared/theme";
import SectionArea from "components/sectionarea";

const options = [
  {
    icon: github,
    label: "hatchet-dev/hatchet",
    value: "hatchet-dev/hatchet",
  },
  {
    icon: github,
    label: "hatchet-dev/hatchet-2",
    value: "hatchet-dev/hatchet-2",
  },
];

const branchOptions = [
  {
    icon: branch,
    label: "master",
    value: "master",
  },
  {
    icon: branch,
    label: "belanger/feat-1",
    value: "belanger/feat-1",
  },
];

const ChooseGitSource: React.FunctionComponent = () => {
  return (
    <SectionArea>
      <H2>Step 1: Choose Git Source</H2>
      <HorizontalSpacer
        spacepixels={14}
        overrides={css({
          borderBottom: theme.line.thick,
        }).toString()}
      />
      <HorizontalSpacer spacepixels={16} />
      <P>
        Choose the Git repository, If you're not seeing the right Git
        repositories, you can link your Git repositories from the Integrations
        tab.
      </P>
      <HorizontalSpacer spacepixels={24} />
      <Selector
        placeholder="Git Repository"
        placeholder_icon={gitRepository}
        options={options}
      />
      <HorizontalSpacer spacepixels={20} />
      <P>
        Choose your branch. This is the branch that your modules should be
        deployed from.
      </P>
      <HorizontalSpacer spacepixels={20} />
      <Selector
        placeholder="Git Branch"
        placeholder_icon={branch}
        options={branchOptions}
      />
      <HorizontalSpacer spacepixels={20} />
      <P>Input the path to the module.</P>
      <HorizontalSpacer spacepixels={20} />
      <TextInput placeholder="path/to/module" />
    </SectionArea>
  );
};

export default ChooseGitSource;
