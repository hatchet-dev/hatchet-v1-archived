import {
  FlexCol,
  FlexRow,
  altTextFontStack,
  FlexRowLeft,
  FlexRowRight,
} from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const StyledSelector = styled(FlexCol)`
  color: ${(props) => props.theme.text.default};
  font-size: 12px;
  background: ${(props) => props.theme.bg.shadeone};
  border: ${(props) => props.theme.line.default};
  padding: 8px 10px;
  border-radius: 5px;
  cursor: pointer;
  width: fit-content;
  :hover {
    background: ${(props) => props.theme.bg.hover};
  }
`;

export const SelectorPlaceholder = styled(FlexRow)`
  > img,
  i:first-child {
    width: 16px;
    height: 16px;
    margin: 0px 8px 0px 0px;
  }

  > div {
    ${altTextFontStack}
    font-size: 13px;
    margin: 0px 8px;
  }

  > i:last-child {
    font-size: 18px;
    padding-top: 1px;
    margin-left: 4px;
  }
`;

export const DropdownWrapper = styled(FlexRow)<{ align: "right" | "left" }>`
  justify-content: ${(props) =>
    props.align == "right" ? "flex-end" : "flex-start"};
  position: absolute;
  width: 100%;
  right: 0;
  z-index: 1;
  top: calc(100% + 5px);
`;

export const Dropdown = styled.div`
  width: fit-content;
  border-radius: 3px;
  z-index: 999;
  overflow-y: auto;
  margin-bottom: 20px;
  background: ${(props) => props.theme.bg.shadeone};
  border-radius: 5px;
  border: ${(props) => props.theme.line.default};
`;

export const ScrollableWrapper = styled.div`
  overflow-y: auto;
  height: 100%;
  max-height: 350px;
`;

export const StyledSelection = styled(FlexRowLeft)`
  padding: 10px 10px;
  cursor: pointer;

  > img,
  i {
    width: 16px;
    height: 16px;
    margin: 0px 8px 0px 0px;
    color: ${(props) => props.theme.text.default};
  }

  > div {
    ${altTextFontStack}
    font-size: 13px;
    margin: 0px 8px;
    color: ${(props) => props.theme.text.default};
  }

  :hover {
    background: ${(props) => props.theme.bg.hover};
  }
`;
