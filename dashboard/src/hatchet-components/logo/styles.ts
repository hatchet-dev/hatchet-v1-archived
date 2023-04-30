import styled from "styled-components";

export const CenteredLogoContainer = styled.div<{
  width?: string;
  height?: string;
  padding?: string;
}>`
  width: ${(props) => props.width || "64px"};
  height: ${(props) => props.height || "64px"};
  border-radius: 8px;
  background: ${(props) => props.theme.bg.buttonone};
  padding: ${(props) => props.padding || "10px"};
`;

export const CenteredLogo = styled.img`
  width: 100%;
  height: 100%;
`;
