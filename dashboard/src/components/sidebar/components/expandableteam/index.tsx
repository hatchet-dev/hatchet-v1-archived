import { useQuery } from "@tanstack/react-query";
import {
  BackText,
  Selector,
  Selection,
  FlexCol,
  H3,
  MaterialIcon,
  FlexRow,
  SmallSpan,
  HorizontalSpacer,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useLocation, useHistory } from "react-router-dom";
import { useAtom } from "jotai";
import { currOrgAtom, currTeamAtom } from "shared/atoms/atoms";

import api from "shared/api";
import { TeamNameAndIcon, TeamExpandIcon, TeamLinkWrapper } from "./styles";
import theme from "shared/theme";
import { css } from "styled-components";
import StyledSidebarLink from "../sidebarlink";
import { SidebarLink } from "components/sidebar";
import { Team } from "shared/api/generated/data-contracts";

type Props = {
  links: SidebarLink[];
  team: Team;
  expanded?: boolean;
};

const ExpandableTeam: React.FunctionComponent<Props> = ({
  links,
  team,
  expanded,
}) => {
  const location = useLocation();
  const history = useHistory();
  const [isExpanded, setIsExpanded] = useState(expanded);
  const [currOrg] = useAtom(currOrgAtom);
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);

  return (
    <FlexCol>
      <TeamNameAndIcon onClick={() => setIsExpanded(!isExpanded)}>
        <SmallSpan>{team.display_name}</SmallSpan>
        <TeamExpandIcon className="material-icons">
          {isExpanded ? "expand_more" : "chevron_right"}
        </TeamExpandIcon>
      </TeamNameAndIcon>
      <TeamLinkWrapper>
        {isExpanded &&
          links.map((val) => {
            const nestedLink = `/teams/${team.id}${val.href}`;
            return (
              <StyledSidebarLink
                key={val.name}
                to={nestedLink}
                current={location?.pathname.includes(nestedLink)}
              >
                {val.name}
              </StyledSidebarLink>
            );
          })}
      </TeamLinkWrapper>
    </FlexCol>
  );
};

export default ExpandableTeam;
