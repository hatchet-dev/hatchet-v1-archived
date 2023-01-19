import { useQuery } from "@tanstack/react-query";
import { BackText, Selector, Selection } from "@hatchet-dev/hatchet-components";
import React from "react";
import { useLocation, useHistory } from "react-router-dom";
import { useAtom } from "jotai";
import { currOrgAtom, currTeamAtom } from "shared/atoms/atoms";

import api from "shared/api";
import {
  LinkWrapper,
  SidebarLink,
  SidebarWrapper,
  UtilWrapper,
} from "./styles";

type Props = {
  links: SidebarLink[];
};

export type SidebarLink = {
  name: string;
  href: string;
};

const SideBar: React.FunctionComponent<Props> = ({ links }) => {
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

  const orgTeamsQuery = useQuery({
    queryKey: ["organization_teams", currOrg.id],
    queryFn: async () => {
      const res = await api.listTeams(currOrg.id);
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

  const onSelectTeam = (option: Selection) => {
    if (option.value == "new_team") {
      history.push("/team/create");
    } else {
      for (let team of orgTeamsQuery.data?.data?.rows) {
        if (option.value == team.id) {
          setCurrTeam(team);
        }
      }
    }
  };

  const renderTeamsUtil = () => {
    if (isUserView) {
      return;
    }

    if (orgTeamsQuery.isLoading) {
      return <div>Loading</div>;
    }

    const selectOptions = orgTeamsQuery.data.data.rows
      .map((row) => {
        return {
          material_icon: "group",
          label: row.display_name,
          value: row.id,
        };
      })
      .concat([
        {
          material_icon: "add_circle",
          label: "New Team",
          value: "new_team",
        },
      ]);

    return (
      <Selector
        placeholder={
          selectOptions.filter((team) => team.value == currTeam?.id)[0]?.label
        }
        placeholder_material_icon="group"
        options={selectOptions}
        select={onSelectTeam}
        orientation="horizontal"
        fill_selection={true}
      />
    );
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

  return (
    <SidebarWrapper>
      <LinkWrapper>
        {renderTeamsUtil()}
        {links.map((val, i) => {
          return (
            <SidebarLink
              key={val.name}
              href={val.href}
              current={location?.pathname.includes(val.href)}
            >
              {val.name}
            </SidebarLink>
          );
        })}
      </LinkWrapper>
      <UtilWrapper>{renderUtil()}</UtilWrapper>
    </SidebarWrapper>
  );
};

export default SideBar;
