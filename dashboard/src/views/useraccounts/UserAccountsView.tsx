import {
  FlexCol,
  FlexColCenter,
  FlexRowRight,
  FlexRowLeft,
  H1,
  H2,
  HorizontalSpacer,
  P,
  StyledDeprecatedText,
  Table,
  StandardButton,
  Spinner,
  Placeholder,
  SmallSpan,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { relativeDate } from "shared/utils";
import { css } from "styled-components";
import theme from "shared/theme";
import GithubAvatarAndName from "components/githubavatarandname";

const UserAccountsView: React.FunctionComponent = () => {
  const [err, setErr] = useState("");

  const metadataQuery = useQuery({
    queryKey: ["api_metadata"],
    queryFn: async () => {
      const res = await api.getServerMetadata();
      return res;
    },
    retry: false,
  });

  const hasGithubAppCapabilities = !!metadataQuery.data?.data?.integrations
    ?.github_app;

  const { data, isLoading } = useQuery({
    queryKey: ["github_app_installations"],
    queryFn: async () => {
      const res = await api.listGithubAppInstallations();
      return res;
    },
    retry: false,
    enabled: hasGithubAppCapabilities,
  });

  const columns = [
    {
      Header: "Account",
      accessor: "account_name",
      Cell: ({ row }: any) => {
        return (
          <GithubAvatarAndName
            account_avatar_url={row.original.account_avatar_url}
            account_name={row.original.account_name}
          />
        );
      },
    },
    {
      Header: "Installed",
      accessor: "created_at",
    },
    {
      Header: "",
      accessor: "id",
      Cell: ({ row }: any) => {
        return (
          <FlexRowRight height="30px">
            <StandardButton
              label="Configure"
              size="small"
              style_kind="muted"
              on_click={async () => {
                window.open(row.original.installation_settings_url);
              }}
            />
          </FlexRowRight>
        );
      },
    },
  ];

  const tableData = data?.data?.rows.map((row) => {
    return {
      id: row.id,
      account_avatar_url: row.account_avatar_url,
      account_name: row.account_name,
      installation_settings_url: row.installation_settings_url,
      created_at: relativeDate(row.created_at),
    };
  });

  const renderGithubAppInstallations = () => {
    if (!hasGithubAppCapabilities) {
      return (
        <Placeholder>
          <SmallSpan>
            This Hatchet instance does not have a Github integration set up.
          </SmallSpan>
        </Placeholder>
      );
    }

    return (
      <>
        <FlexRowRight>
          <StandardButton
            label="Install"
            material_icon="add"
            on_click={() => {
              window.open("/api/v1/github_app/install");
            }}
          />
        </FlexRowRight>
        <HorizontalSpacer spacepixels={20} />
        <Table
          columns={columns}
          data={tableData}
          dataName="personal access tokens"
        />
      </>
    );
  };

  if (metadataQuery.isLoading || (hasGithubAppCapabilities && isLoading)) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  return (
    <FlexColCenter>
      <FlexCol width="100%" maxWidth="640px">
        <H1>Linked Accounts</H1>

        <HorizontalSpacer
          spacepixels={10}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={10} />
        <H2>Github App Installations</H2>
        <HorizontalSpacer spacepixels={10} />
        <P>Manage the accounts that the Hatchet Github app has access to.</P>
        <HorizontalSpacer spacepixels={10} />
        {renderGithubAppInstallations()}
      </FlexCol>
    </FlexColCenter>
  );
};

export default UserAccountsView;
