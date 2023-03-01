import {
  FlexCol,
  FlexColCenter,
  FlexRowRight,
  H1,
  H2,
  HorizontalSpacer,
  P,
  StyledDeprecatedText,
  Table,
  StandardButton,
  Spinner,
  Placeholder,
  ErrorBar,
  SectionArea,
  TextInput,
} from "@hatchet-dev/hatchet-components";
import React, { useState, useEffect, useCallback } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { relativeDate } from "shared/utils";
import { useAtom } from "jotai";
import { currOrgAtom } from "shared/atoms/atoms";
import { useHistory } from "react-router-dom";
import { CreateTeamRequest, Team } from "shared/api/generated/data-contracts";
import TeamWithMembers from "components/team/teamwithmembers";

const CreateTeam: React.FunctionComponent = () => {
  const history = useHistory();
  const [currOrg, setCurrOrg] = useAtom(currOrgAtom);
  const [displayName, setDisplayName] = useState("");
  const [createdTeam, setCreatedTeam] = useState<Team>(null);
  const [err, setErr] = useState("");

  const mutation = useMutation({
    mutationKey: ["create_team", currOrg.id],
    mutationFn: async (team: CreateTeamRequest) => {
      const res = await api.createTeam(currOrg.id, team);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
      setCreatedTeam(data?.data);
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  const handleKeyPress = useCallback(
    (e: any) => {
      e.key === "Enter" && submit();
    },
    [displayName]
  );

  useEffect(() => {
    window.addEventListener("keydown", handleKeyPress);

    return () => {
      window.removeEventListener("keydown", handleKeyPress);
    };
  }, [handleKeyPress]);

  const submit = () => {
    setErr("");

    if (displayName != "") {
      mutation.mutate({
        display_name: displayName,
      });
    }
  };

  const renderInnerForm = () => {
    if (createdTeam) {
      return (
        <>
          <SectionArea>
            <TeamWithMembers team={createdTeam} collapsible={false} />
          </SectionArea>
          <HorizontalSpacer spacepixels={20} />
          <FlexRowRight>
            <StandardButton
              label="Finish"
              material_icon="chevron_right"
              icon_side="right"
              on_click={() => {
                history.push("/organization/teams");
              }}
              margin={"0"}
            />
          </FlexRowRight>
        </>
      );
    }

    return (
      <>
        <SectionArea>
          <TextInput
            placeholder="ex. Team 1"
            label="Team Name"
            type="text"
            width="400px"
            on_change={(val) => {
              setDisplayName(val);
            }}
          />
          <HorizontalSpacer spacepixels={20} />
          {err && <ErrorBar text={err} />}
        </SectionArea>
        <HorizontalSpacer spacepixels={20} />
        <FlexRowRight>
          <StandardButton
            label="Create"
            material_icon="chevron_right"
            icon_side="right"
            on_click={() => {
              submit();
            }}
            disabled={displayName == ""}
            margin={"0"}
            is_loading={mutation.isLoading}
          />
        </FlexRowRight>
      </>
    );
  };

  return (
    <FlexColCenter>
      <FlexCol width="100%" maxWidth="640px">
        <H2>Create a new team</H2>
        <HorizontalSpacer spacepixels={20} />
        {renderInnerForm()}
      </FlexCol>
    </FlexColCenter>
  );
};

export default CreateTeam;
