import { FlexCol, FlexRow, Span } from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const TopBarWrapper = styled.nav`
  border-bottom: ${theme.line.default};
  background-color: ${theme.bg.default};
  width: 100%;
  opacity: 1;
  padding: 16px 24px;
  height: 70px;
  position: fixed;
  top: 0;
  z-index: 2;
`;

export const TopBarContainer = styled(FlexCol)``;

export const TopBarProductName = styled(Span)`
  font-weight: 700;
  font-family: "Ubuntu", sans-serif;
  font-size: 20px;
  letter-spacing: 1px;
`;

export const ProfileContainer = styled(FlexRow)`
  max-width: fit-content;
  margin-right: 10px;
  border-radius: 6px;
  padding: 4px 6px;
  cursor: pointer;

  :hover {
    background-color: ${theme.bg.hover};
  }
`;

export const ProfileName = styled(Span)`
  padding: 0 6px;
`;
