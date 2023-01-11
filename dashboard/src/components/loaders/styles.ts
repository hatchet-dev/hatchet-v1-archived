import { MaterialIcon } from "components/globals";
import styled from "styled-components";

export const LoadingSpinner = styled(MaterialIcon)`
  animation: spinspinspin 1.2s linear infinite;
  color: #ffffffaa;

  @keyframes spinspinspin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
`;
