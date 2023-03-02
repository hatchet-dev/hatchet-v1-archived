import {
  SmallSpan,
  FlexRowLeft,
  SectionCard,
  FlexRowRight,
} from "@hatchet-dev/hatchet-components";
import { StatusContainer } from "components/status/styles";
import styled from "styled-components";

export const TriggerPRContainer = styled(StatusContainer)`
  font-size: 12px;
  padding: 4px 8px;
  cursor: pointer;

  > i {
    font-size: 12px;
  }
`;
