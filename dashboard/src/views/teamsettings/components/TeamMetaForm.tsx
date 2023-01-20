import {
  FlexRowRight,
  HorizontalSpacer,
  TextInput,
  SectionArea,
  StandardButton,
  ErrorBar,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { currOrgAtom, currTeamAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import { TeamUpdateRequest } from "shared/api/generated/data-contracts";

const TeamMetaForm: React.FunctionComponent = () => {
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);
  const [displayName, setDisplayName] = useState("");
  const [err, setErr] = useState("");

  const { mutate, isLoading } = useMutation({
    mutationKey: ["update_team", currTeam.id],
    mutationFn: (teamUpdate: TeamUpdateRequest) => {
      return api.updateTeam(currTeam.id, teamUpdate);
    },
    onSuccess: (data) => {
      if (data?.data) {
        setCurrTeam(data?.data);
      }
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const submit = () => {
    if (displayName != "" && displayName != currTeam?.display_name) {
      mutate({
        display_name: displayName,
      });
    }
  };

  return (
    <SectionArea>
      <TextInput
        placeholder="ex. Team 1"
        initial_value={currTeam?.display_name}
        label="Display name"
        type="text"
        width="400px"
        on_change={(val) => {
          setDisplayName(val);
        }}
      />
      <HorizontalSpacer spacepixels={30} />
      {err && <ErrorBar text={err} />}
      <HorizontalSpacer spacepixels={30} />
      <FlexRowRight>
        <StandardButton
          label="Update"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            submit();
          }}
          disabled={displayName == "" || displayName == currTeam?.display_name}
          margin={"0"}
          is_loading={isLoading}
        />
      </FlexRowRight>
    </SectionArea>
  );
};

export default TeamMetaForm;
