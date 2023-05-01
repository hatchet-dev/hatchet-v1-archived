import { FlexCol, Span } from "../globals";
import styled from "styled-components";

export const TopBarWrapper = styled.nav<{ is_authenticated?: boolean }>`
  ${(props) =>
    props.is_authenticated && `border-bottom: ${props.theme.line.default};`}

  background-color: ${(props) => props.theme.bg.default};
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
  margin-left: 12px;
`;
