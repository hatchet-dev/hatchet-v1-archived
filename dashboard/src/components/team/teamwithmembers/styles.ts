import {
  altTextFontStack,
  FlexColScroll,
  FlexRow,
  MaterialIcon,
} from "@hatchet-dev/hatchet-components";
import styled from "styled-components";

export const TeamHeader = styled(FlexRow)`
  gap: 12px;
  margin: 4px 0;
  color: ${(props) => props.theme.text.default};
`;

export const TeamNameWithIcon = styled(FlexRow)`
  gap: 10px;
`;

export const ExpandIcon = styled(MaterialIcon)`
  color: ${(props) => props.theme.text.default};
  cursor: pointer;
`;
