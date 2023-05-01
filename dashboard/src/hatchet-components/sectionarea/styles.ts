import { FlexCol } from "../globals";
import styled from "styled-components";

export const StyledSectionArea = styled(FlexCol)<{
  width?: string;
  height?: string;
  background?: string;
  border?: string;
  flex?: string;
}>`
  border: ${(props) =>
    props.border ? props.border : props.theme.line.default};
  border-radius: 8px;
  padding: 32px;
  width: ${(props) => (props.width ? props.width : "100%")};
  height: ${(props) => (props.height ? props.height : "auto")};
  background: ${(props) =>
    props.background ? props.background : props.theme.bg.default};
  ${(props) => props.flex && `flex: ${props.flex};`};
`;
