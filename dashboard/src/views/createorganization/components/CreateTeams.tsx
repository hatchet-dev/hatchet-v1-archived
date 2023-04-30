import {
  FlexCol,
  FlexRowRight,
  HorizontalSpacer,
  StandardButton,
  SectionArea,
} from "hatchet-components";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import { CreateTeamRequest } from "shared/api/generated/data-contracts";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import TeamManager from "components/team/teammanager";

const teamHelper =
  "Add teams by entering the team name and assigning team members to each team. You can also add teams later from organization settings.";

const CreateTeams: React.FunctionComponent = () => {
  const [currOrg] = useAtom(currOrgAtom);
  const [err, setErr] = useState("");
  const [createdTeam, setCreatedTeam] = useState(false);
  const history = useHistory();

  const mutation = useMutation({
    mutationKey: ["create_team", currOrg.id],
    mutationFn: async (team: CreateTeamRequest) => {
      const res = await api.createTeam(currOrg.id, team);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  if (!currOrg) {
    history.push("/");
  }

  return (
    <FlexCol>
      <SectionArea width="600px">
        <TeamManager
          add_team_helper={teamHelper}
          create_team={(t) => {
            setCreatedTeam(true);
          }}
        />
      </SectionArea>
      <HorizontalSpacer spacepixels={24} />
      <FlexRowRight>
        <StandardButton
          label="Next"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            history.push("/");
          }}
          margin={"0"}
          is_loading={mutation.isLoading}
          disabled={!createdTeam}
        />
      </FlexRowRight>
    </FlexCol>
  );
};

export default CreateTeams;
