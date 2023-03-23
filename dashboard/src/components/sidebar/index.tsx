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
import React from "react";
import { useLocation, useHistory } from "react-router-dom";
import { useAtom } from "jotai";
import { currOrgAtom, currTeamAtom } from "shared/atoms/atoms";

import api from "shared/api";
import {
  LinkWrapper,
  SidebarWrapper,
  TeamName,
  TeamExpandIcon,
  UtilWrapper,
  TeamNameHeader,
} from "./styles";
import theme from "shared/theme";
import { css } from "styled-components";
import StyledSidebarLink from "./components/sidebarlink";
import ExpandableTeam from "./components/expandableteam";

type Props = {
  links: SidebarLink[];
  team_links?: SidebarLink[];
};

export type SidebarLink = {
  name: string;
  href: string;
};

const SideBar: React.FunctionComponent<Props> = ({
  links,
  team_links = [],
}) => {
  const location = useLocation();
  const history = useHistory();
  const isUserView = location.pathname.includes("/user");
  const [currOrg, setCurrOrg] = useAtom(currOrgAtom);
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);

  const { data, isLoading } = useQuery({
    queryKey: ["current_user_organizations"],
    queryFn: async () => {
      const res = await api.listUserOrganizations();
      return res;
    },
    retry: false,
  });

  const userTeamsQuery = useQuery({
    queryKey: ["current_user_teams", currOrg?.id],
    queryFn: async () => {
      const res = await api.listUserTeams({
        organization_id: currOrg.id,
      });
      return res;
    },
    retry: false,
    enabled: !!currOrg,
  });

  const onSelectOrg = (option: Selection) => {
    if (option.value == "new_organization") {
      history.push("/organization/create");
    } else {
      for (let org of data?.data?.rows) {
        if (option.value == org.id) {
          setCurrOrg(org);
        }
      }
    }
  };

  const renderUtil = () => {
    if (isUserView) {
      return (
        <BackText
          text="Dashboard"
          back={() => {
            history.push("/");
          }}
          width="100%"
        />
      );
    }

    if (isLoading) {
      return <div>Loading</div>;
    }

    const selectOptions = data.data.rows
      .map((row) => {
        return {
          material_icon: "corporate_fare",
          label: row.display_name,
          value: row.id,
        };
      })
      .concat([
        {
          material_icon: "add_circle",
          label: "New Organization",
          value: "new_organization",
        },
      ]);

    return (
      <Selector
        placeholder={
          selectOptions.filter((org) => org.value == currOrg.id)[0]?.label
        }
        placeholder_material_icon="corporate_fare"
        options={selectOptions}
        select={onSelectOrg}
        orientation="vertical"
        option_alignment="right"
        fill_selection={true}
      />
    );
  };

  const renderSidebarLinks = () => {
    if (isUserView) {
      return (
        <LinkWrapper>
          {links.map((val, i) => {
            return (
              <StyledSidebarLink
                key={val.name}
                to={val.href}
                current={location?.pathname.includes(val.href)}
              >
                {val.name}
              </StyledSidebarLink>
            );
          })}
        </LinkWrapper>
      );
    }

    const userTeams = userTeamsQuery?.data?.data?.rows;

    return (
      <LinkWrapper>
        {links.map((val, i) => {
          return (
            <StyledSidebarLink
              key={val.name}
              to={val.href}
              current={location?.pathname.includes(val.href)}
            >
              {val.name}
            </StyledSidebarLink>
          );
        })}
        <HorizontalSpacer spacepixels={10} />
        <TeamNameHeader>Your teams</TeamNameHeader>
        <HorizontalSpacer
          spacepixels={6}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={3} />
        {userTeams?.map((row, i) => {
          return (
            <>
              <ExpandableTeam
                team={row}
                links={team_links}
                expanded={currTeam && currTeam.id == row.id}
              />
              {userTeams.length != i + 1 && (
                <HorizontalSpacer spacepixels={10} />
              )}
            </>
          );
        })}
      </LinkWrapper>
    );
  };

  return (
    <SidebarWrapper>
      {renderSidebarLinks()}
      <UtilWrapper>{renderUtil()}</UtilWrapper>
    </SidebarWrapper>
  );
};

export default SideBar;
