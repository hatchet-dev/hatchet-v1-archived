import { FlexRowLeft } from "@hatchet-dev/hatchet-components";
import styled from "styled-components";

export const CheckboxContainer = styled(FlexRowLeft)`
  > i {
    color: ${(props) => props.theme.text.default};
    cursor: pointer;
  }
`;
