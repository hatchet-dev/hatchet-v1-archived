import {
  H2,
  HorizontalSpacer,
  P,
  Selector,
  TextInput,
  SectionArea,
  Placeholder,
  Spinner,
  FlexRowRight,
  StandardButton,
  Breadcrumbs,
  H1,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import gitRepository from "assets/git_repository.png";
import github from "assets/github.png";
import branch from "assets/branch.png";
import { css } from "styled-components";
import theme from "shared/theme";
import {
  CreateModuleRequest,
  GithubAppInstallation,
  GithubBranch,
  GithubRepo,
} from "shared/api/generated/data-contracts";
import { useQuery } from "@tanstack/react-query";
import api from "shared/api";
import usePrevious from "shared/hooks/useprevious";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";

type Props = {
  submit: (req: CreateModuleRequest) => void;
};

const ChooseGitSource: React.FunctionComponent<Props> = ({ submit }) => {
  const [selectedGAI, setSelectedGAI] = useState<GithubAppInstallation>(null);
  const [selectedGAIReset, setSelectedGAIReset] = useState(0);
  const [selectedRepo, setSelectedRepo] = useState<GithubRepo>(null);
  const [selectedRepoReset, setSelectedRepoReset] = useState(0);
  const [selectedBranch, setSelectedBranch] = useState<GithubBranch>(null);
  const [path, setPath] = useState("");
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);

  const breadcrumbs = [
    {
      label: "Modules",
      link: `/team/${currTeam.id}/modules`,
    },
    {
      label: "Step 1: Choose Git Source",
      link: `/team/${currTeam.id}/modules/create/step_1`,
    },
  ];

  const previousGAI: GithubAppInstallation = usePrevious(selectedGAI);
  const previousRepo: GithubRepo = usePrevious(selectedRepo);

  useEffect(() => {
    if (previousGAI && previousGAI.id != selectedGAI?.id) {
      setSelectedGAIReset(selectedGAIReset + 1);
    }

    if (previousRepo && previousRepo.repo_name != selectedRepo?.repo_name) {
      setSelectedRepoReset(selectedRepoReset + 1);
    }
  }, [selectedGAI, previousGAI, selectedRepo, previousRepo]);

  const gaiQuery = useQuery({
    queryKey: ["github_app_installations"],
    queryFn: async () => {
      const res = await api.listGithubAppInstallations();
      return res;
    },
    retry: false,
  });

  const reposQuery = useQuery({
    queryKey: ["github_repos", selectedGAI?.id],
    queryFn: async () => {
      const res = await api.listGithubRepos(selectedGAI?.id);
      return res;
    },
    retry: false,
    enabled: !!selectedGAI,
  });

  const branchesQuery = useQuery({
    queryKey: [
      "github_branches",
      selectedGAI?.id,
      selectedRepo?.repo_owner,
      selectedRepo?.repo_name,
    ],
    queryFn: async () => {
      const res = await api.listGithubRepoBranches(
        selectedGAI?.id,
        selectedRepo?.repo_owner,
        selectedRepo?.repo_name
      );
      return res;
    },
    retry: false,
    enabled: !!selectedGAI && !!selectedRepo,
  });

  if (gaiQuery.isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  const gais = gaiQuery.data?.data?.rows;

  const githubAccountOptions =
    gais?.map((gai) => {
      return {
        icon: gai.account_avatar_url,
        label: gai.account_name,
        value: gai.id,
      };
    }) || [];

  const repos = reposQuery.data?.data;

  const repoOptions =
    repos?.map((repo) => {
      return {
        icon: github,
        label: `${repo.repo_owner}/${repo.repo_name}`,
        value: repo.repo_name,
      };
    }) || [];

  const branches = branchesQuery.data?.data;

  const branchOptions =
    branches?.map((branch) => {
      return {
        icon: github,
        label: branch.branch_name,
        value: branch.branch_name,
      };
    }) || [];

  const submitEnabled =
    !!selectedGAI && !!selectedRepo && !!selectedBranch && path != "";

  return (
    <>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
      <HorizontalSpacer spacepixels={12} />
      <H1>Create a new module</H1>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H2>Step 1: Choose Git Source</H2>
        <HorizontalSpacer
          spacepixels={14}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={16} />
        <P>Choose the Github account to deploy this module from.</P>
        <HorizontalSpacer spacepixels={24} />
        <Selector
          placeholder="Github Account"
          placeholder_icon={gitRepository}
          options={githubAccountOptions}
          select={(selection) => {
            setSelectedGAI(
              gais.filter((gai) => {
                return gai.id == selection.value;
              })[0]
            );

            setSelectedRepo(null);
            setSelectedBranch(null);
            setPath("");
          }}
        />
        <HorizontalSpacer spacepixels={24} />
        {selectedGAI && (
          <>
            <P>Choose the Github repo to deploy this module from.</P>
            <HorizontalSpacer spacepixels={24} />
            <Selector
              placeholder="Github Repo"
              placeholder_icon={gitRepository}
              options={repoOptions}
              reset={selectedGAIReset}
              select={(selection) => {
                setSelectedRepo(
                  repos.filter((repo) => {
                    return repo.repo_name == selection.value;
                  })[0]
                );

                setSelectedBranch(null);
                setPath("");
              }}
            />
          </>
        )}
        <HorizontalSpacer spacepixels={24} />
        {selectedRepo && (
          <>
            <P>Choose the Github branch to deploy this module from.</P>
            <HorizontalSpacer spacepixels={24} />
            <Selector
              placeholder="Github Branch"
              placeholder_icon={gitRepository}
              options={branchOptions}
              reset={selectedRepoReset}
              select={(selection) => {
                setSelectedBranch(
                  branches.filter((branch) => {
                    return branch.branch_name == selection.value;
                  })[0]
                );

                setPath("");
              }}
            />
          </>
        )}
        {selectedBranch && (
          <>
            <HorizontalSpacer spacepixels={20} />
            <P>Input the path to the module.</P>
            <HorizontalSpacer spacepixels={20} />
            <TextInput
              placeholder="path/to/module"
              on_change={(val) => setPath(val)}
            />
          </>
        )}
      </SectionArea>
      <HorizontalSpacer spacepixels={20} />
      <FlexRowRight>
        <StandardButton
          label="Next"
          material_icon="chevron_right"
          icon_side="right"
          disabled={!submitEnabled}
          on_click={() => {
            if (!submitEnabled) {
              return;
            }

            submit({
              name: `${selectedRepo.repo_owner}-${selectedRepo.repo_name}-${selectedBranch.branch_name}`,
              github: {
                path: path,
                github_app_installation_id: selectedGAI.id,
                github_repository_owner: selectedRepo.repo_owner,
                github_repository_name: selectedRepo.repo_name,
                github_repository_branch: selectedBranch.branch_name,
              },
            });
          }}
        />
      </FlexRowRight>
    </>
  );
};

export default ChooseGitSource;
