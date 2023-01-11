import { FlexCol } from "components/globals";
import React, { useEffect, useState } from "react";
import usePrevious from "shared/hooks/useprevious";
import { SmallSpan, StyledInput } from "./styles";

export type Selection = {
  label: string;
  value: string;
  icon?: string;
};

export type Props = {
  placeholder: string;
  initial_value?: string;
  label?: string;
  type?: "text" | "password";
  width?: string;
  disabled?: boolean;
  on_change?: (val: string) => void;
  reset?: number;
};

const TextInput: React.FC<Props> = ({
  placeholder,
  initial_value = "",
  label,
  type = "text",
  width = "250px",
  disabled = false,
  on_change,
  reset,
}) => {
  const [text, setText] = useState<string>(initial_value);
  const prevReset = usePrevious(reset);

  useEffect(() => {
    if (reset != prevReset) {
      setText("");
    }
  }, [reset]);

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
      disabled={disabled}
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
