import { FlexColScroll, FlexRow } from "hatchet-components";
import styled from "styled-components";

export const OrgListContainer = styled(FlexColScroll)`
  max-height: 400px;
`;

export const OrgContainer = styled(FlexRow)`
  gap: 12px;
  margin: 6px 0;
`;

export const OrgName = styled.div`
  font-size: 13px;
  width: 100%;
  padding: 7px 12px;
  background: ${(props) => props.theme.bg.shadeone};
  border: ${(props) => props.theme.line.default};
  border-radius: 4px;
  color: ${(props) => props.theme.text.default};
`;
