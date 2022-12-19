import { FlexRow, tableHeaderFontStack } from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const TabWrapper = styled(FlexRow)`
  justify-content: start;
  border-bottom: ${theme.line.default};
  margin: 10px 0;
`;

export const Tab = styled.div<{ selected?: boolean }>`
  ${tableHeaderFontStack};
  color: ${theme.text.default};
  font-weight: ${(props) => (props.selected ? "600" : "400")};
  font-size: 13px;
  cursor: pointer;
  padding: 4px 6px 8px 6px;
  border-bottom: ${(props) =>
    props.selected ? `${theme.line.selected}` : "none"};
  border-bottom-width: ${(props) => (props.selected ? "2px" : "0px")};
  margin: 0 10px;
  margin-bottom: ${(props) => (props.selected ? "-2px" : "0")};

  :first-child {
    margin-left: 0px;
  }

  :last-child {
    margin-right: 0px;
  }
`;
