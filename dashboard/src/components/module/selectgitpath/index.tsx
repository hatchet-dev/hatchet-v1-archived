import {
  HorizontalSpacer,
  P,
  Selector,
  TextInput,
  Placeholder,
  Spinner,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useState } from "react";
import gitRepository from "assets/git_repository.png";
import github from "assets/github.png";
import {
  CreateModuleRequestGithub,
  CreateModuleValuesRequestGithub,
  GithubAppInstallation,
  GithubBranch,
  GithubRepo,
} from "shared/api/generated/data-contracts";
import { useQuery } from "@tanstack/react-query";
import api from "shared/api";
import usePrevious from "shared/hooks/useprevious";

type Props = {
  // current_params can be passed in from current state (when updating a module), or
  // if query data has previously been stored (not implemented)
  current_params?: CreateModuleRequestGithub;
  set_request: (req: CreateModuleRequestGithub) => void;
};

const SelectGitSource: React.FunctionComponent<Props> = ({
  set_request,
  current_params,
}) => {
  const [selectedGAI, setSelectedGAI] = useState<GithubAppInstallation>(null);
  const [selectedGAIReset, setSelectedGAIReset] = useState(0);
  const [selectedRepo, setSelectedRepo] = useState<GithubRepo>(null);
  const [selectedRepoReset, setSelectedRepoReset] = useState(0);
  const [selectedBranch, setSelectedBranch] = useState<GithubBranch>(null);
  const [path, setPath] = useState(current_params?.path || "");

  const previousGAI: GithubAppInstallation = usePrevious(selectedGAI);
  const previousRepo: GithubRepo = usePrevious(selectedRepo);

  useEffect(() => {
    set_request({
      github_app_installation_id: selectedGAI?.id || "",
      github_repository_owner: selectedGAI?.account_name || "",
      github_repository_name: selectedRepo?.repo_name || "",
      github_repository_branch: selectedBranch?.branch_name || "",
      path: path || "",
    });
  }, [selectedGAI, selectedRepo, selectedBranch, path]);

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

      // if no GAI is selected and there is a current id, set that as the selected GAI
      if (!selectedGAI && current_params?.github_app_installation_id) {
        const matchedGAI = res.data?.rows?.filter((gai) => {
          return gai.id == current_params?.github_app_installation_id;
        })[0];

        matchedGAI && setSelectedGAI(matchedGAI);
      }

      return res;
    },
    retry: false,
  });

  const reposQuery = useQuery({
    queryKey: ["github_repos", selectedGAI?.id],
    queryFn: async () => {
      const res = await api.listGithubRepos(selectedGAI?.id);

      // if no repo is selected and there is a current repo, set that as the selected repo
      if (!selectedRepo && current_params?.github_repository_name) {
        const matchedRepo = res.data?.filter((repo) => {
          return repo.repo_name == current_params?.github_repository_name;
        })[0];

        matchedRepo && setSelectedRepo(matchedRepo);
      }

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

      // if no branch is selected and there is a current branch, set that as the selected branch
      if (!selectedBranch && current_params?.github_repository_branch) {
        const matchedBranch = res.data?.filter((branch) => {
          return branch.branch_name == current_params?.github_repository_branch;
        })[0];

        matchedBranch && setSelectedBranch(matchedBranch);
      }

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
    branches
      ?.sort((a) => {
        return a.is_default ? -1 : 1;
      })
      .map((branch) => {
        return {
          icon: github,
          label: branch.branch_name,
          value: branch.branch_name,
        };
      }) || [];

  return (
    <>
      <P>Choose the Github account to deploy this module from.</P>
      <HorizontalSpacer spacepixels={24} />
      <Selector
        placeholder={selectedGAI?.account_name || "Github Account"}
        placeholder_icon={selectedGAI?.account_avatar_url || gitRepository}
        options={githubAccountOptions}
        select={(selection) => {
          const selectedGAI = gais.filter((gai) => {
            return gai.id == selection.value;
          })[0];

          setSelectedGAI(selectedGAI);
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
            placeholder={selectedRepo?.repo_name || "Github Repo"}
            placeholder_icon={selectedRepo ? github : gitRepository}
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
            placeholder={selectedBranch?.branch_name || "Github Branch"}
            placeholder_icon={selectedBranch ? github : gitRepository}
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
            initial_value={path}
            on_change={(val) => {
              setPath(val);
            }}
          />
        </>
      )}
    </>
  );
};

export default SelectGitSource;
