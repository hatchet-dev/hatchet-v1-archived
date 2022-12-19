import { FlexRow, textFontStack } from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const FilterWrapper = styled(FlexRow)`
  border: 1px solid black;
  padding: 1px 14px;
  max-width: fit-content;
  border-radius: 6px;
  margin: 8px 0;
`;

export const FilterText = styled.div`
  ${textFontStack}
  color: ${theme.text.default};
  margin-left: 8px;
  font-size: 12px;
`;
