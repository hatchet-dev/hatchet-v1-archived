import theme from "shared/theme";
import styled from "styled-components";

export const StyledInput = styled.input`
  outline: none;
  border: none;
  font-size: 13px;
  width: 250px;
  padding: 10px 12px;
  background: ${theme.bg.shadeone};
  border: ${theme.line.default};
  border-radius: 4px;
  color: ${theme.text.default};

  ::placeholder {
    color: ${theme.text.inactive};
  }
`;
