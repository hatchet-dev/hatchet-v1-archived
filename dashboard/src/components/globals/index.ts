import { Link } from "react-router-dom";
import styled, { css } from "styled-components";

export const Shadow = {
  low: "0 2px 8px",
  mid: "0 4px 12px",
  high: "0 8px 16px",
};

export const Transition = {
  hover: {
    on: "all 0.2s ease-in",
    off: "all 0.2s ease-out",
  },
  reaction: {
    on: "all 0.15s ease-in",
    off: "all 0.1s ease-out",
  },
  dropdown: {
    off: "all 0.35s ease-out",
  },
};

export const headerFontStack = css`
  font-family: "${(props) => props.theme.font.header}", -apple-system,
    BlinkMacSystemFont, "Helvetica", "Segoe", sans-serif;
`;

export const tableHeaderFontStack = css`
  font-family: "Work Sans", Arial, sans-serif;
  font-size: 12px;
  font-weight: 600;
`;

export const sidebarHeaderFontStack = css`
  font-family: "Work Sans", Arial, sans-serif;
`;

export const altTextFontStack = css`
  font-family: "Work Sans", Arial, sans-serif;
`;

export const textFontStack = css`
  font-family: "${(props) => props.theme.font.text}", Arial, sans-serif;
`;

export const monoStack = css`
  font-family: "Input Mono", "Menlo", "Inconsolata", "Roboto Mono", monospace;
`;

export const Label = styled.label`
  display: flex;
  flex-direction: column;
  width: 100%;
  margin-top: 0.5rem;
  font-weight: 500;
  font-size: 0.875rem;
  letter-spacing: -0.4px;
  color: ${(props) => props.theme.text.default};
  &:not(:first-of-type) {
    margin-top: 1.5rem;
  }
  a {
    text-decoration: underline;
  }
`;

export const H1 = styled.h1`
  ${headerFontStack};
  color: ${(props) => props.theme.text.default};
  font-weight: 900;
  font-size: 1.8rem;
  line-height: 1.3;
  margin: 0;
  padding: 0;
`;

export const H2 = styled.h2`
  color: ${(props) => props.theme.text.default};
  ${headerFontStack};
  font-weight: 700;
  font-size: 1.25rem;
  line-height: 1.3;
  margin: 0;
  padding: 0;
`;

export const H3 = styled.h3`
  color: ${(props) => props.theme.text.default};
  ${headerFontStack};
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.5;
  margin: 0;
  padding: 0;
`;

export const H4 = styled.h4`
  color: ${(props) => props.theme.text.default};
  ${headerFontStack};
  font-weight: 500;
  font-size: 16px;
  line-height: 1.4;
  margin: 0;
  padding: 0;
`;

export const H5 = styled.h5`
  color: ${(props) => props.theme.text.default};
  ${headerFontStack};
  font-weight: 500;
  font-size: 0.75rem;
  line-height: 1.4;
  margin: 0;
  padding: 0;
`;

export const H6 = styled.h6`
  color: ${(props) => props.theme.text.default};
  ${headerFontStack};
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.675rem;
  line-height: 1.5;
  margin: 0;
  padding: 0;
`;

export const P = styled.p`
  color: ${(props) => props.theme.text.default};
  ${textFontStack};
  font-weight: 400;
  font-size: 0.875rem;
  line-height: 1.6;
  margin: 0;
  padding: 0;
`;

export const Span = styled.span`
  color: ${(props) => props.theme.text.default};
  ${textFontStack};
  font-weight: 400;
  font-size: 0.875rem;
  line-height: 1.6;
  margin: 0;
  padding: 0;
`;

export const SmallSpan = styled(Span)`
  font-size: 0.8rem;
`;

export const StyledClickableP = styled(P)`
  color: ${(props) => props.theme.text.link};
  font-weight: 400;
  font-size: 0.875rem;
  line-height: 1.4;
  cursor: pointer;
`;

export const StyledLink = styled(Link)`
  color: ${(props) => props.theme.text.link};
  text-decoration: none;
  cursor: pointer;
  font-weight: 400;
  font-size: 0.875rem;
  line-height: 1.4;
  margin: 0;
  padding: 0;
`;

export const StyledSmallLink = styled(StyledLink)`
  font-size: 0.8rem;
`;

export const Relative = styled.div`
  position: relative;
`;

export const FlexRow = styled.div<{ height?: string; maxHeight?: string }>`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  ${(props) => props.height && `height: ${props.height};`}
  ${(props) => props.maxHeight && `max-height: ${props.maxHeight};`}
`;

export const FlexRowRight = styled(FlexRow)`
  justify-content: flex-end;
`;

export const FlexRowLeft = styled(FlexRow)`
  justify-content: flex-start;
`;

export const FlexCol = styled.div<{
  width?: string;
  maxWidth?: string;
  height?: string;
  maxHeight?: string;
}>`
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: stretch;
  ${(props) => props.width && `width: ${props.width};`}
  ${(props) => props.maxWidth && `max-width: ${props.maxWidth};`}
  ${(props) => props.height && `height: ${props.height};`}
  ${(props) => props.maxHeight && `max-height: ${props.maxHeight};`}
`;

export const FlexColCenter = styled(FlexCol)`
  align-items: center;
`;

export const FlexColScroll = styled(FlexCol)`
  overflow-x: auto;
`;

export const Grid = styled.div`
  display: grid;
  grid-column-gap: 25px;
  grid-row-gap: 25px;
  grid-template-columns: 140px 140px 140px;
  grid-template-rows: 100px auto;
`;

export const MaterialIcon = styled.i`
  font-size: 16px;
  margin: 6px 0;
`;

export const HoverableMaterialIcon = styled(MaterialIcon)`
  cursor: pointer;
  font-size: 18px;
  border-radius: 4px;

  :hover {
    background-color: ${(props) => props.theme.bg.hover};
  }
`;

export const HorizontalSpacer = styled.div<{
  spacepixels: number;
  overrides?: string;
}>`
  width: 100%;
  border-top: 0px solid red;
  margin: ${(props) => Math.round(props.spacepixels / 2.0)}px 0;
  ${(props) => props?.overrides}
`;

export const fadeInAnimation = css`
  animation: fadeIn 0.5s;
  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
`;

export const StyledDeprecatedText = styled(FlexRow)`
  ${tableHeaderFontStack};
  padding: 8px 16px;
  background-color: ${(props) => props.theme.bg.shadetwo};
  border: ${(props) => props.theme.line.default};
  color: white;
  border-radius: 4px;
  height: 28px;
  font-weight: 500;
  width: fit-content;
`;
