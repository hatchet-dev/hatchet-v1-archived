import { createGlobalStyle } from "styled-components";

export const theme = {
  bg: {
    default: "#04061c",
    reverse: "#16171A",
    wash: "#FAFAFA",
    divider: "white",
    border: "white",
    inactive: "#DFE7EF",
    shadeone: "#141321",
    shadetwo: "#211F2D",
    buttonone: "#3F16E4",
    hover: "#211F2D",
    hovertwo: "#ffffff0b",
  },
  font: {
    header: "Ubuntu",
    text: "Noto Sans Tai Le",
  },
  line: {
    default: "1px solid #ffffff22",
    thick: "2px solid #ffffff22",
    selected: "1px solid #ffffffaa",
  },
  text: {
    default: "#ffffffcc",
    alt: "white",
    inactive: "#999999aa",
  },
};

export default theme;

export const GlobalStyle = createGlobalStyle`
* {
  box-sizing: border-box;
  font-family: 'Noto Sans Tai Le', sans-serif;
}

body {
  background: ${theme.bg.default};
  overscroll-behavior-x: none;
}

a {
  text-decoration: none;
}

img {
  max-width: 100%;
}
`;
