import { FlexRow, textFontStack } from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const PaginatorWrapper = styled(FlexRow)`
  padding: 1px 14px;
  max-width: fit-content;
  border-radius: 6px;
  margin: 8px 0;
`;

export const PaginatorText = styled.div`
  ${textFontStack}
  color: ${theme.text.default};
  font-size: 12px;
  margin: 0 8px;
`;
