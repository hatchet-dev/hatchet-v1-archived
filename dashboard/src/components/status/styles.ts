import {
  FlexRowLeft,
  altTextFontStack,
  SmallSpan,
} from "@hatchet-dev/hatchet-components";
import styled from "styled-components";

export const StatusContainer = styled(FlexRowLeft)<{ width?: string }>`
  ${altTextFontStack}
  color: ${(props) => props.theme.text.default};
  font-size: 13px;
  border: ${(props) => props.theme.line.default};
  padding: 2px 8px;
  border-radius: 5px;
  cursor: default;
  width: ${(props) => props.width || "fit-content"};
  :hover {
    background: ${(props) => props.theme.bg.hover};
  }

  > i {
    color: ${(props) => props.theme.text.default};
    font-size: 16px;
    margin-right: 10px;
  }
`;

export const StatusText = styled(SmallSpan)`
  font-weight: bold;
`;
