import {
  altTextFontStack,
  FlexColScroll,
  FlexRow,
} from "@hatchet-dev/hatchet-components";
import styled from "styled-components";

export const MemberListContainer = styled(FlexColScroll)`
  max-height: 400px;
`;

export const MemberContainer = styled(FlexRow)`
  gap: 12px;
  margin: 4px 0;
`;

export const MemberNameOrEmail = styled.div`
  font-size: 13px;
  width: 400px;
  padding: 7px 12px;
  background: ${(props) => props.theme.bg.shadeone};
  border: ${(props) => props.theme.line.default};
  border-radius: 4px;
  color: ${(props) => props.theme.text.default};
`;

export const PolicyName = styled(FlexRow)`
  color: ${(props) => props.theme.text.default};
  font-size: 12px;
  background: ${(props) => props.theme.bg.shadeone};
  border: ${(props) => props.theme.line.default};
  padding: 8px 10px;
  border-radius: 5px;
  width: 120px;
  flex-grow: 0;
  flex-shrink: 0;

  > img,
  i:first-child {
    width: 16px;
    height: 16px;
    margin: 0px 8px 0px 0px;
  }

  > div {
    ${altTextFontStack}
    font-size: 13px;
    margin: 0px 8px;
  }
`;
