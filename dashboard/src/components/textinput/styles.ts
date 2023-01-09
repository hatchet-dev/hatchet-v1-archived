import { Span } from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const StyledInput = styled.input<{ width?: string; disabled?: boolean }>`
  outline: none;
  border: none;
  font-size: 13px;
  width: ${(props) => (props.width ? props.width : "250px")};
  padding: 8px 12px;
  background: ${(props) => props.theme.bg.shadeone};
  border: ${(props) => props.theme.line.default};
  border-radius: 4px;
  color: ${(props) =>
    props.disabled ? props.theme.text.inactive : props.theme.text.default};

  ::placeholder {
    color: ${(props) => props.theme.text.inactive};
  }
`;

export const SmallSpan = styled(Span)`
  font-size: 13px;
  font-weight: 700;
  margin-bottom: 8px;
`;
