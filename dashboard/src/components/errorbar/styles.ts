import {
  fadeInAnimation,
  FlexRowLeft,
  SmallSpan,
  textFontStack,
} from "components/globals";
import styled from "styled-components";

export const StyledErrorBar = styled(FlexRowLeft)<{ color?: string }>`
  ${fadeInAnimation}
  border: ${(props) => props.theme.line.default};
  ${(props) => props.color && `border-color: ${props.color};`}
  background: ${(props) => props.theme.bg.shadeone};
  color: ${(props) => props.color || props.theme.text.default};
  margin-bottom: 3px;
  border-radius: 6px;
  padding: 8px 14px;
  overflow: hidden;
  gap: 12px;
`;

export const StyledErrorText = styled(SmallSpan)<{ color?: string }>`
  ${(props) => props.color && `color: ${props.color};`}
`;
