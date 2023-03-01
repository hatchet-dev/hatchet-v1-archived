import {
  FlexRowLeft,
  MaterialIcon,
  SmallSpan,
  P,
} from "@hatchet-dev/hatchet-components";
import React, { Component, useRef, useState } from "react";
import { CheckboxContainer } from "./styles";

type Props = {
  text: string;
  setChecked: (v: boolean) => void;
  initialValue?: boolean;
};

const Checkbox: React.FC<Props> = ({
  text,
  setChecked,
  initialValue = false,
}) => {
  const [isChecked, setIsChecked] = useState(initialValue);

  return (
    <CheckboxContainer gap="8px">
      <MaterialIcon
        className="material-icons"
        onClick={() => {
          setChecked(!isChecked);
          setIsChecked(!isChecked);
        }}
      >
        {isChecked ? "check_box" : "check_box_outline_blank"}
      </MaterialIcon>
      <P>{text}</P>
    </CheckboxContainer>
  );
};

export default Checkbox;
