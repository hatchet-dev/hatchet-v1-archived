import {
  fadeInAnimation,
  FlexRowLeft,
  MaterialIcon,
  SmallSpan,
  textFontStack,
} from "components/globals";
import styled from "styled-components";

export const LoadingSpinner = styled(MaterialIcon)`
  animation: lds-dual-ring 1.2s linear infinite;

  @keyframes lds-dual-ring {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
`;

// export const LoadingSpinner = styled.div`
//   width: 80px;
//   height: 80px;

//   :after {
//     content: " ";
//     display: block;
//     width: 64px;
//     height: 64px;
//     margin: 8px;
//     border-radius: 50%;
//     border: 6px solid #fff;
//     border-color: #fff transparent #fff transparent;
//     animation: lds-dual-ring 1.2s linear infinite;
//   }

//   @keyframes lds-dual-ring {
//     0% {
//       transform: rotate(0deg);
//     }
//     100% {
//       transform: rotate(360deg);
//     }
//   }
// `;
