import { FlexRowLeft, altTextFontStack } from "../globals";
import styled from "styled-components";

export const BackArrow = styled(FlexRowLeft)<{ width?: string }>`
  ${altTextFontStack}
  color: ${(props) => props.theme.text.default};
  font-size: 13px;
  background: ${(props) => props.theme.bg.shadeone};
  border: ${(props) => props.theme.line.default};
  padding: 8px 20px 8px 10px;
  border-radius: 5px;
  cursor: pointer;
  width: ${(props) => props.width || "fit-content"};
  :hover {
    background: ${(props) => props.theme.bg.hover};
  }

  > i {
    color: ${(props) => props.theme.text.default};
    font-size: 18px;
    margin-right: 6px;
  }
`;
