import { FlexCol } from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const StyledSectionArea = styled(FlexCol)<{
  width?: number;
  height?: number;
  background?: string;
}>`
  border: ${(props) => props.theme.line.default};
  border-radius: 8px;
  padding: 32px;
  width: ${(props) => (props.width ? props.width + "px" : "100%")};
  height: ${(props) => (props.height ? props.height + "px" : "auto")};
  background: ${(props) =>
    props.background ? props.background : props.theme.bg.default};
`;
