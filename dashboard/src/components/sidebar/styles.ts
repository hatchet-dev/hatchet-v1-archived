import {
  FlexCol,
  FlexRow,
  sidebarHeaderFontStack,
} from "@hatchet-dev/hatchet-components";
import styled from "styled-components";

export const SidebarWrapper = styled.section`
  border-right: ${(props) => props.theme.line.default};
  background-color: ${(props) => props.theme.bg.default};
  height: 100%;
  width: 230px;
  opacity: 1;
  padding: 8px;
  padding-top: 60px;
  position: fixed;
  left: 0;
  z-index: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: stretch;
`;

export const LinkWrapper = styled(FlexCol)`
  padding: 4px 3px;
  margin: 20px 4px;
`;

export const SidebarLink = styled.a<{ current: boolean }>`
  ${sidebarHeaderFontStack}
  cursor: pointer;
  font-size: 13px;
  margin: 4px 0;
  padding: 10px;
  text-decoration: none;
  color: ${(props) =>
    props.current ? props.theme.text.alt : props.theme.text.default};
  border-radius: 6px;

  :hover {
    background-color: ${(props) => props.theme.bg.hover};
  }
`;

export const UtilWrapper = styled(FlexRow)`
  height: fit-content;
`;
