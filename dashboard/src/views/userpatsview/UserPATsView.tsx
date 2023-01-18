import {
  FlexCol,
  FlexColCenter,
  FlexRowRight,
  H1,
  HorizontalSpacer,
  P,
  StyledDeprecatedText,
  Table,
  StandardButton,
  Spinner,
  Placeholder,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import CreatePATForm from "./components/CreatePATForm";
import { relativeDate } from "shared/utils";

const UserPATsView: React.FunctionComponent = () => {
  const [createPAT, setCreatePAT] = useState(false);
  const [err, setErr] = useState("");

  const { data, isLoading, refetch, isFetching } = useQuery({
    queryKey: ["pats"],
    queryFn: async () => {
      const res = await api.listPersonalAccessTokens();
      return res;
    },
    retry: false,
  });

  const revokeMutation = useMutation({
    mutationKey: ["revoke_pat"],
    mutationFn: (id: string) => {
      return api.revokePersonalAccessToken(id);
    },
    onSuccess: (data) => {
      setErr("");
      refetch();
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const columns = [
    {
      Header: "Name",
      accessor: "name",
    },
    {
      Header: "Created",
      accessor: "created_at",
    },
    {
      Header: "",
      accessor: "id",
      Cell: ({ row }: any) => {
        if (row.original.revoked) {
          return (
            <FlexRowRight height="30px">
              <StyledDeprecatedText>Revoked</StyledDeprecatedText>
            </FlexRowRight>
          );
        }
        return (
          <FlexRowRight height="30px">
            <StandardButton
              label="Revoke"
              size="small"
              style_kind="muted"
              disabled={row.original.revoked}
              on_click={async () => {
                await revokeMutation.mutateAsync(row.original.id);

                refetch();
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
      name: row.display_name,
      created_at: relativeDate(row.created_at),
      revoked: row.revoked,
    };
  });

  if (isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  const renderPATDataOrForm = () => {
    if (createPAT) {
      return (
        <CreatePATForm
          post_create={() => {
            refetch();
            setCreatePAT(false);
          }}
        />
      );
    }

    return (
      <>
        <FlexRowRight>
          <StandardButton
            label="Create new PAT"
            material_icon="add"
            on_click={() => {
              setCreatePAT(true);
            }}
          />
        </FlexRowRight>
        <HorizontalSpacer spacepixels={20} />
        <Table
          columns={columns}
          data={tableData}
          dataName="personal access tokens"
        />
        <HorizontalSpacer spacepixels={20} />
      </>
    );
  };

  return (
    <FlexColCenter>
      <FlexCol width="100%" maxWidth="640px">
        <H1>Personal Access Tokens</H1>
        <HorizontalSpacer spacepixels={12} />
        <P>
          Personal access tokens can be used to authenticate with the Hatchet
          API. Personal access tokens are automatically generated for new CLI
          sessions.
        </P>
        <HorizontalSpacer spacepixels={16} />
        {renderPATDataOrForm()}
      </FlexCol>
    </FlexColCenter>
  );
};

export default UserPATsView;
