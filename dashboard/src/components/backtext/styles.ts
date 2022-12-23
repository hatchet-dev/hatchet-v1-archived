import theme from "shared/theme";
import styled from "styled-components";

export const BackArrowContainer = styled.div`
  width: 100%;
  height: 24px;
`;

export const BackArrow = styled.div`
  > i {
    color: ${(props) => props.theme.text.default};
    font-size: 18px;
    margin-right: 6px;
  }

  color: ${(props) => props.theme.text.default};
  display: flex;
  align-items: center;
  font-size: 14px;
  cursor: pointer;
  width: 120px;
`;
