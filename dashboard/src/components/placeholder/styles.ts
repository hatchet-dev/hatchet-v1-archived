import { FlexColCenter } from "components/globals";
import styled from "styled-components";

export const StyledPlaceholder = styled(FlexColCenter)<{
  width?: string;
  height?: string;
  background?: string;
}>`
  border-radius: 8px;
  padding: 32px;
  width: ${(props) => (props.width ? props.width : "100%")};
  height: ${(props) => (props.height ? props.height : "auto")};
  background: ${(props) =>
    props.background ? props.background : props.theme.bg.shadetwo};
`;
