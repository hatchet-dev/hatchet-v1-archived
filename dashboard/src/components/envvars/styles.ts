import { FlexRowLeft, FlexRow } from "@hatchet-dev/hatchet-components";
import styled from "styled-components";

export const EnvVarRow = styled(FlexRowLeft)`
  gap: 8px;
`;

export const EnvVarRemoveButton = styled(FlexRow)<{
  has_icon?: boolean;
  disabled?: boolean;
}>`
  padding: 2px 8px;
  background-color: ${(props) => props.theme.bg.shadetwo};
  color: white;
  border-radius: 50%;
  cursor: pointer;

  > i {
    font-weight: 700;
    font-size: 14px;
  }

  :hover {
    filter: brightness(93%);
  }
`;
