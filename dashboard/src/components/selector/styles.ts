import {
  FlexCol,
  FlexRow,
  altTextFontStack,
  FlexRowLeft,
  Relative,
} from "components/globals";
import styled from "styled-components";

export const StyledSelectorWrapper = styled(Relative)<{
  orientation: "horizontal" | "vertical";
}>`
  ${(props) => (props.orientation == "vertical" ? "width: 100%;" : "")}
  overflow: visible;
`;

export const StyledSelector = styled(FlexCol)<{
  orientation: "horizontal" | "vertical";
}>`
  color: ${(props) => props.theme.text.default};
  font-size: 12px;
  background: ${(props) => props.theme.bg.shadeone};
  border: ${(props) => props.theme.line.default};
  padding: 8px 10px;
  border-radius: 5px;
  cursor: pointer;
  width: ${(props) =>
    props.orientation == "vertical" ? "100%" : "fit-content"};
  overflow-x: visible;
  :hover {
    background: ${(props) => props.theme.bg.hover};
  }
`;

export const SelectorPlaceholder = styled(FlexRow)`
  > i:last-child {
    font-size: 18px;
    padding-top: 1px;
    margin-left: 4px;
  }
`;

export const InnerSelectorPlaceholder = styled(FlexRowLeft)`
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
`;

export const DropdownWrapper = styled(FlexRow)<{
  align: "right" | "left";
  orientation: "horizontal" | "vertical";
}>`
  justify-content: ${(props) =>
    props.align == "right" ? "flex-end" : "flex-start"};
  position: absolute;
  ${(props) =>
    props.orientation == "vertical"
      ? `
  left: calc(100% + 5px);
  bottom: -20px; 
  `
      : ""}
  ${(props) =>
    props.orientation == "horizontal"
      ? `
  top: 42px;
  `
      : ""}
  ${(props) =>
    props.orientation == "horizontal" &&
    (props.align == "right"
      ? `
  right: 0;
  `
      : "left: 0;")}
  z-index: 5;
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
  white-space: nowrap;

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
