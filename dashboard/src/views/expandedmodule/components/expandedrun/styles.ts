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

export const ExpandedRunContainer = styled.div`
  overflow-x: auto;
  height: calc(100% - 200px);
`;

export const RunSectionCard = styled(SectionCard)`
  padding: 20px;
  height: fit-content;
`;
