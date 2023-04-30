import {
  altTextFontStack,
  FlexRow,
  SmallSpan,
  Span,
  P,
  tableHeaderFontStack,
  FlexCol,
} from "hatchet-components";
import styled from "styled-components";

export const ExpandableSettingsContainer = styled(FlexCol)`
  border-radius: 6px;
  border: ${(props) => props.theme.line.default};
`;

export const ExpandableSettingsHeader = styled(FlexRow)`
  background: ${(props) => props.theme.bg.shadetwo};
  cursor: pointer;
  padding: 4px 14px;

  > i {
    color: ${(props) => props.theme.text.default};
  }
`;

export const ExpandableSettingsBody = styled(FlexCol)<{ hidden?: boolean }>`
  padding: 20px;
  ${(props) => props.hidden && "display: none;"}
`;

export const ExpandableSettingsText = styled(SmallSpan)`
  font-weight: bold;
`;
