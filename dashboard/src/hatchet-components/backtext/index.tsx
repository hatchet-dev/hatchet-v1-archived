import React from "react";
import { BackArrow } from "./styles";

export type Props = {
  text: string;
  back: () => void;
  width?: string;
};

const BackText: React.FC<Props> = ({ text, back, width }) => {
  return (
    <BackArrow onClick={() => back()} width={width}>
      <i className="material-icons next-icon">navigate_before</i>
      {text}
    </BackArrow>
  );
};

export default BackText;
