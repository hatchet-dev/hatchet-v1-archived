import { FlexRowRight, tableHeaderFontStack } from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const StyledStandardButton = styled(FlexRowRight)<{
  size?: "small" | "default";
  has_icon?: boolean;
  icon_side?: "left" | "right";
  margin?: string;
  disabled?: boolean;
}>`
  ${tableHeaderFontStack};
  padding: ${(props) =>
    props.size == "small"
      ? "6px 10px"
      : props.has_icon
      ? "2px 16px"
      : "8px 16px"};
  padding-right: ${(props) =>
    !props.has_icon ? "" : props.icon_side == "left" ? "20px" : "16px"};
  padding-left: ${(props) =>
    !props.has_icon ? "" : props.icon_side == "right" ? "20px" : "16px"};
  background-color: ${(props) =>
    props.disabled ? props.theme.bg.inactive : props.theme.bg.buttonone};
  color: white;
  border-radius: 4px;
  font-size: 13px;
  margin: ${(props) => (props.margin ? props.margin : "0 8px")};
  cursor: ${(props) => (props.disabled ? "default" : "pointer")};
  font-weight: ${(props) => (props.size == "small" ? "400" : "500")};
  gap: 8px;

  > i {
    font-weight: 700;
  }

  :hover {
    filter: brightness(93%);
  }
`;
