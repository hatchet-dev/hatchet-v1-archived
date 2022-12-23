import { FlexCol, Relative, Span } from "components/globals";
import React, { useEffect, useRef, useState } from "react";
import { SmallSpan, StyledInput } from "./styles";

export type Selection = {
  label: string;
  value: string;
  icon?: string;
};

export type Props = {
  placeholder: string;
  label?: string;
  type?: "text" | "password";
  width?: string;
  on_change?: (val: string) => void;
};

const TextInput: React.FC<Props> = ({
  placeholder,
  label,
  type = "text",
  width = "250px",
  on_change,
}) => {
  const [text, setText] = useState<string>("");

  const input = (
    <StyledInput
      placeholder={placeholder}
      type={type}
      value={text}
      width={width}
      onChange={(e) => {
        setText(e.target.value);
        on_change && on_change(e.target.value);
      }}
    />
  );

  if (label) {
    return (
      <FlexCol>
        <SmallSpan>{label}</SmallSpan>
        {input}
      </FlexCol>
    );
  }

  return input;
};

export default TextInput;
