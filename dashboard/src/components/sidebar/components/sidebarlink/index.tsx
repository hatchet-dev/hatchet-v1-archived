import { sidebarHeaderFontStack } from "@hatchet-dev/hatchet-components";
import { Link } from "react-router-dom";
import styled from "styled-components";

const StyledSidebarLink = styled(Link)<{ current: boolean; padding?: string }>`
  ${sidebarHeaderFontStack}
  cursor: pointer;
  font-size: 13px;
  margin: 4px 0;
  padding: ${(props) => props.padding || "10px"};
  text-decoration: none;
  color: ${(props) =>
    props.current ? props.theme.text.alt : props.theme.text.default};
  border-radius: 6px;

  :hover {
    background-color: ${(props) => props.theme.bg.hover};
  }
`;

export default StyledSidebarLink;
