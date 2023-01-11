import { FlexRowRight, tableHeaderFontStack } from "components/globals";
import styled from "styled-components";

export const StyledStandardButton = styled(FlexRowRight)<{
  size?: "small" | "default";
  has_icon?: boolean;
  icon_side?: "left" | "right";
  style_kind?: "default" | "muted";
  margin?: string;
  disabled?: boolean;
}>`
  ${tableHeaderFontStack};
  padding: ${(props) =>
    props.size == "small"
      ? "6px 14px"
      : props.has_icon
      ? "2px 16px"
      : "8px 16px"};
  padding-right: ${(props) =>
    !props.has_icon ? "" : props.icon_side == "left" ? "20px" : "16px"};
  padding-left: ${(props) =>
    !props.has_icon ? "" : props.icon_side == "right" ? "20px" : "16px"};
  background-color: ${(props) =>
    props.disabled
      ? props.theme.bg.inactive
      : props.style_kind == "muted"
      ? props.theme.bg.shadeone
      : props.theme.bg.buttonone};
  ${(props) =>
    props.style_kind == "muted" && `border: ${props.theme.line.default};`}
  color: white;
  border-radius: 4px;
  font-size: ${(props) => (props.size == "small" ? "12px" : "13px")};
  ${(props) => props.size == "small" && "height: 28px;"}
  margin: ${(props) => (props.margin ? props.margin : "0 8px")};
  cursor: ${(props) => (props.disabled ? "default" : "pointer")};
  font-weight: 500;
  gap: 8px;
  white-space: nowrap;

  > i {
    font-weight: 700;
  }

  :hover {
    filter: brightness(93%);
  }
`;
