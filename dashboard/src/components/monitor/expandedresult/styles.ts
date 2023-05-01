import { SectionCard } from "hatchet-components";
import styled from "styled-components";

export const ExpandedResultContainer = styled.div`
  overflow-x: auto;
  height: calc(100% - 200px);
`;

export const ResultSectionCard = styled(SectionCard)`
  padding: 20px;
  height: fit-content;
`;
