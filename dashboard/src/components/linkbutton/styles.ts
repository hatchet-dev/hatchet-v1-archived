import { StatusContainerStyles } from "components/status/styles";
import { Link } from "react-router-dom";
import styled from "styled-components";

export const LinkButtonContainer = styled(Link)`
  ${StatusContainerStyles}
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;

  > i {
    margin-right: 0;
    margin-left: 10px;
  }
`;
