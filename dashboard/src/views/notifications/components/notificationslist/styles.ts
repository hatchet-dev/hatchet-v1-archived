import {
  FlexCol,
  FlexColScroll,
  SmallSpan,
} from "hatchet-components";
import styled from "styled-components";

export const NotificationListContainer = styled(FlexColScroll)`
  border: ${(props) => props.theme.line.default};
  height: 100%;
  border-radius: 4px;
  padding-right: 0;
`;

export const NotificationMetaContainer = styled(FlexCol)`
  padding: 16px 12px;
  cursor: pointer;
  border-bottom: ${(props) => props.theme.line.thick};

  :hover {
    background: ${(props) => props.theme.bg.shadetwo};
  }
`;

export const NotificationTitle = styled(SmallSpan)`
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
`;

export const NotificationDate = styled(SmallSpan)`
  white-space: nowrap;
`;
