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

export const SidebarWrapper = styled.section<{ padding?: string }>`
  border-right: ${(props) => props.theme.line.default};
  background-color: ${(props) => props.theme.bg.default};
  height: 100%;
  width: 230px;
  opacity: 1;
  padding: ${(props) => props.padding || "60px 2px 8px 2px"};
  position: fixed;
  left: 0;
  z-index: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: stretch;
`;

export const LinkWrapper = styled(FlexColScroll)<{ padding?: string }>`
  padding: ${(props) => props.padding || "4px 3px"};
  margin: 20px 4px;
  max-height: calc(100% - 100px);
`;

export const UtilWrapper = styled(FlexRow)`
  height: fit-content;
`;

export const TeamName = styled(Span)`
  ${headerFontStack}
  font-weight: bold;
`;

export const TeamExpandIcon = styled(MaterialIcon)`
  color: ${(props) => props.theme.text.default};
  cursor: pointer;
`;

export const TeamNameHeader = styled(SmallSpan)`
  ${tableHeaderFontStack}
  font-weight: bold;
  padding-left: 10px;
  color: ${(props) => props.theme.text.inactive};
`;
