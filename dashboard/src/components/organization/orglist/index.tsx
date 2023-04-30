import {
  HorizontalSpacer,
  ErrorBar,
  StandardButton,
} from "hatchet-components";
import React from "react";
import { Organization } from "shared/api/generated/data-contracts";
import { OrgListContainer, OrgContainer, OrgName } from "./styles";

export type Props = {
  orgs: Organization[];
  leave_org?: (org: Organization) => void;
  err?: string;
};

const OrgList: React.FC<Props> = ({ orgs, leave_org, err }) => {
  return (
    <OrgListContainer>
      {orgs.map((org, i) => {
        return (
          <OrgContainer key={org.id}>
            <OrgName>{org.display_name}</OrgName>
            <StandardButton
              label="Leave Organization"
              style_kind="muted"
              on_click={() => leave_org(org)}
            />
          </OrgContainer>
        );
      })}
      {err && <HorizontalSpacer spacepixels={20} />}
      {err && <ErrorBar text={err} />}
    </OrgListContainer>
  );
};

export default OrgList;
