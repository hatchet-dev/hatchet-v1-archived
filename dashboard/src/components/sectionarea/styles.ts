import { FlexCol } from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const StyledSectionArea = styled(FlexCol)<{ width?: number }>`
  border: ${theme.line.default};
  border-radius: 8px;
  padding: 32px;
  max-height: 600px;
  width: ${(props) => (props.width ? props.width + "px" : "100%")};
`;
