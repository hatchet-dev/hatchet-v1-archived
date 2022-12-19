import React from "react";
import { BackArrowContainer, BackArrow } from "./styles";

export type Props = {
  text: string;
  back: () => void;
};

const BackText: React.FC<Props> = ({ text, back }) => {
  return (
    <BackArrowContainer>
      <BackArrow onClick={() => back()}>
        <i className="material-icons next-icon">navigate_before</i>
        {text}
      </BackArrow>
    </BackArrowContainer>
  );
};

export default BackText;
