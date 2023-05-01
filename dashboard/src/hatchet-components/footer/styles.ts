import styled from "styled-components";

export const FooterWrapper = styled.footer`
  background-color: ${(props) => props.theme.bg.default};
  border-top: ${(props) => props.theme.line.default};
  width: 100%;
  opacity: 1;
  padding: 16px 24px;
  height: 100px;
  z-index: 2;
`;
