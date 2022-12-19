import React, { useState } from "react";
import { useLocation } from "react-router-dom";
import { LinkWrapper, SidebarLink, SidebarWrapper } from "./styles";

type SidebarLink = {
  name: string;
  href: string;
};

const SidebarLinks: SidebarLink[] = [
  {
    name: "Home",
    href: "/home",
  },
  {
    name: "Modules",
    href: "/modules",
  },
  {
    name: "Monitoring",
    href: "/monitoring",
  },
  {
    name: "Templates",
    href: "/templates",
  },
  {
    name: "Integrations",
    href: "/integrations",
  },
  {
    name: "Settings",
    href: "/settings",
  },
];

const SideBar: React.FunctionComponent = () => {
  let location = useLocation();

  console.log(location);

  return (
    <SidebarWrapper>
      <LinkWrapper>
        {SidebarLinks.map((val, i) => {
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
    </SidebarWrapper>
  );
};

export default SideBar;
