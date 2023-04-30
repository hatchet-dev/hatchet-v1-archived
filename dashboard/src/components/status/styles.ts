import {
  FlexRowLeft,
  altTextFontStack,
  SmallSpan,
} from "hatchet-components";
import styled, { css } from "styled-components";

export const StatusContainerStyles = css`
  color: ${(props) => props.theme.text.default};
  font-size: 13px;
  border: ${(props) => props.theme.line.default};
  padding: 2px 8px;
  border-radius: 5px;
  cursor: default;

  :hover {
    background: ${(props) => props.theme.bg.hover};
  }

  > i {
    color: ${(props) => props.theme.text.default};
    font-size: 16px;
    margin-right: 10px;
  }
`;

export const StatusContainer = styled(FlexRowLeft)<{
  width?: string;
  color?: string;
}>`
  ${altTextFontStack}
  ${StatusContainerStyles}
  width: ${(props) => props.width || "fit-content"};
  ${(props) => props.color && `border-color: ${props.color};`}
`;

export const StatusText = styled(SmallSpan)<{ color?: string }>`
  font-weight: bold;

  ${(props) => props.color && `color: ${props.color};`}
`;
