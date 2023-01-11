import { FlexCol } from "components/globals";
import styled from "styled-components";

export const StyledSectionArea = styled(FlexCol)<{
  width?: string;
  height?: string;
  background?: string;
}>`
  border: ${(props) => props.theme.line.default};
  border-radius: 8px;
  padding: 32px;
  width: ${(props) => (props.width ? props.width : "100%")};
  height: ${(props) => (props.height ? props.height : "auto")};
  background: ${(props) =>
    props.background ? props.background : props.theme.bg.default};
`;
