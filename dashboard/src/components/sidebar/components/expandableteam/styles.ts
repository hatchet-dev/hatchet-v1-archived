import {
  FlexCol,
  FlexColScroll,
  FlexRow,
  sidebarHeaderFontStack,
  Span,
  MaterialIcon,
  tableHeaderFontStack,
  headerFontStack,
  SmallSpan,
} from "hatchet-components";
import styled from "styled-components";

export const TeamLinkWrapper = styled(FlexCol)`
  padding-left: 10px;
`;

export const TeamNameAndIcon = styled(FlexRow)`
  cursor: pointer;
  padding: 2px 10px;
  border-radius: 6px;

  :hover {
    background-color: ${(props) => props.theme.bg.hover};
  }
`;

export const TeamExpandIcon = styled(MaterialIcon)`
  color: ${(props) => props.theme.text.default};
  cursor: pointer;
`;
