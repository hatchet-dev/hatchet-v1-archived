import StandardButton from "components/buttons";
import ErrorBar from "components/errorbar";
import {
  FlexCol,
  FlexRow,
  HorizontalSpacer,
  MaterialIcon,
} from "components/globals";
import SectionArea from "components/sectionarea";
import Selector from "components/selector";
import TextInput from "components/textinput";
import React, { useEffect, useState } from "react";
import {
  Organization,
  OrganizationMemberSanitized,
} from "shared/api/generated/data-contracts";
import { capitalize } from "shared/utils";
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
