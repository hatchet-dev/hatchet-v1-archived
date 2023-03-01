import { FlexCol, FlexColScroll } from "@hatchet-dev/hatchet-components";
import styled from "styled-components";

export const ExpandedNotificationContainer = styled(FlexCol)`
  height: 100%;
  width: 100%;
  background: ${(props) => props.theme.bg.shadeone};
  border-radius: 0 12px 12px 0;
  padding: 20px;
`;
