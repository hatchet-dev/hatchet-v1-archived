import React, { useState } from "react";
import { useLocation } from "react-router-dom";
import { LinkWrapper, SidebarLink, SidebarWrapper } from "./styles";

type Props = {
  links: SidebarLink[];
};

export type SidebarLink = {
  name: string;
  href: string;
};

const SideBar: React.FunctionComponent<Props> = ({ links }) => {
  let location = useLocation();

  console.log(location);

  return (
    <SidebarWrapper>
      <LinkWrapper>
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
    </SidebarWrapper>
  );
};

export default SideBar;
