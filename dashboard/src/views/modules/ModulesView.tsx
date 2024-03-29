import {
  FlexRowRight,
  H1,
  HorizontalSpacer,
  P,
  StandardButton,
  Filter,
  Paginator,
  Table,
  Spinner,
  Placeholder,
} from "hatchet-components";
import { useQuery } from "@tanstack/react-query";
import GithubAvatarAndName from "components/githubavatarandname";
import { useAtom } from "jotai";
import React from "react";
import { useHistory } from "react-router-dom";
import api from "shared/api";
import { currTeamAtom } from "shared/atoms/atoms";
import { relativeDate } from "shared/utils";
import github from "assets/github.png";
import usePagination from "shared/hooks/usepagination";

const ModulesView: React.FunctionComponent = () => {
  let history = useHistory();
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);

  const {
    currentPage,
    maxPage,
    cursor_forward,
    cursor_backward,
    set_data,
  } = usePagination();

  const listModulesQuery = useQuery({
    queryKey: ["modules", currTeam?.id, currentPage],
    queryFn: async () => {
      const res = await api.listModules(currTeam?.id, {
        page: currentPage,
      });

      return res;
    },
    retry: false,
    onSuccess: (data) => {
      set_data(data?.data?.pagination);
    },
  });

  const columns = [
    {
      Header: "Name",
      accessor: "name",
      width: 200,
    },
    {
      Header: "Deployment",
      accessor: "deployment_mechanism",
      width: 150,
      Cell: ({ row }: any) => {
        if (row.original.deployment_mechanism == "local") {
          return <div>local</div>;
        }

        return (
          <GithubAvatarAndName
            account_avatar_url={github}
            account_name={row.original.repo_name}
            avatar_size="small"
          />
        );
      },
    },
    {
      Header: "Last Deployed",
      accessor: "last_deployed",
      width: 100,
    },
    {
      Header: "Module Path",
      accessor: "path",
    },
  ];

  const handleResourceClick = (row: any) => {
    history.push(`/teams/${currTeam.id}/modules/${row.original.id}`);
  };

  const handleCreateModuleClick = () => {
    history.push(`/teams/${currTeam.id}/modules/create/step_1`);
  };

  const moduleData = listModulesQuery.data?.data.rows.map((module) => {
    return {
      id: module.id,
      name: module.name,
      deployment_mechanism: module.deployment_mechanism,
      repo_name: `${module.deployment.git_repo_owner}/${module.deployment.git_repo_name}`,
      repo_branch: module.deployment.git_repo_branch,
      last_deployed: relativeDate(module.updated_at),
      path: module.deployment.path,
    };
  });

  const renderModules = () => {
    if (listModulesQuery.isLoading) {
      return (
        <Placeholder>
          <Spinner />
        </Placeholder>
      );
    }

    return (
      <Table
        columns={columns}
        data={moduleData}
        onRowClick={handleResourceClick}
      />
    );
  };

  return (
    <>
      <H1>Modules</H1>
      <HorizontalSpacer spacepixels={12} />
      <P>Modules are all Terraform modules which have a Terraform state.</P>
      <FlexRowRight>
        <StandardButton
          label="Create module"
          material_icon="add"
          on_click={handleCreateModuleClick}
          margin="0"
        />
      </FlexRowRight>
      <HorizontalSpacer spacepixels={30} />
      {renderModules()}
      <FlexRowRight>
        <Paginator
          curr_page={currentPage}
          last_page={maxPage}
          cursor_backward={cursor_backward}
          cursor_forward={cursor_forward}
        />
      </FlexRowRight>
    </>
  );
};

export default ModulesView;
