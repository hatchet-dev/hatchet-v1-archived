import styled from "styled-components";

const AppWrapper = styled.div<{ background?: string }>`
  width: 100vw;
  height: 100vh;
  position: fixed;
  top: 0;
  left: 0;
  margin: 0;
  display: flex;
  justify-content: center;
  background: ${(props) => props.theme.bg.default};
`;

export default AppWrapper;
