import {
  altTextFontStack,
  FlexColScroll,
  FlexRow,
} from "hatchet-components";
import styled from "styled-components";

export const TeamContainer = styled(FlexRow)`
  gap: 12px;
  margin: 4px 0;
`;

export const TeamName = styled.div`
  font-size: 13px;
  width: 100%;
  padding: 7px 12px;
  background: ${(props) => props.theme.bg.shadeone};
  border: ${(props) => props.theme.line.default};
  border-radius: 4px;
  color: ${(props) => props.theme.text.default};
`;
