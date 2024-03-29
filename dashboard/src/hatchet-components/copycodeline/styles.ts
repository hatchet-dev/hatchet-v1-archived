import styled from "styled-components";

export const CodeBlock = styled.div<{ padding?: string }>`
  display: inline-block;
  color: ${(props) => props.theme.text.codehighlight};
  border-radius: 5px;
  font-family: monospace;
  user-select: text;
  overflow: auto;
  padding: ${(props) => props.padding || "10px"};
  white-space: nowrap;
  border: ${(props) => props.theme.line.default};
`;
