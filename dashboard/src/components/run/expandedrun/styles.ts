import {
  SmallSpan,
  FlexRowLeft,
  SectionCard,
  FlexRowRight,
} from "@hatchet-dev/hatchet-components";
import { StatusContainer } from "components/status/styles";
import styled from "styled-components";

export const ExpandedRunContainer = styled.div`
  overflow-x: auto;
  height: calc(100% - 200px);
`;

export const RunSectionCard = styled(SectionCard)`
  padding: 20px;
  height: fit-content;
`;
